package mw

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"wraith.me/message_server/crud"
	"wraith.me/message_server/db"
	"wraith.me/message_server/obj"
	cr "wraith.me/message_server/redis"
	"wraith.me/message_server/util/httpu"
)

var (
	//The name of the cookie to look for.
	AuthCookieName = "token"

	//The name of the HTTP query parameter to look for.
	AuthHttpParamName = "token"

	//The name of the auth subject header to send
	AuthHttpHeaderSubject = "X-Auth-For"

	//The name of the auth scope header to send
	AuthHttpHeaderScope = "X-Auth-Scope"
)

// Holds the error messages.
var (
	ErrAuthUnauthorized   = errors.New("token is not authorized for this route")
	ErrAuthExpiredToken   = errors.New("token has expired")
	ErrAuthNoTokenFound   = errors.New("no token found")
	ErrAuthBadTokenFormat = errors.New("token format is incorrect")
	ErrAuthGeneric        = errors.New("authentication error")
)

type authMiddleware struct {
	allowedScopes []obj.TokenScope //The scopes for which the token is valid, sorted in increasing order.
	mclient       *mongo.Client    //The MongoDB database client.
	rclient       *redis.Client    //The Redis database client.
}

// Returns a new handler for the authentication middleware.
func NewAuthMiddleware(allowedScopes []obj.TokenScope) func(next http.Handler) http.Handler {
	//Get a struct object
	mw := authMiddleware{
		allowedScopes: allowedScopes,
		mclient:       db.GetInstance().GetClient(),
		rclient:       cr.GetInstance().GetClient(),
	}

	//Return the instance
	slices.Sort(mw.allowedScopes)
	return mw.authMWHandler
}

/*
TokenFromCookie tries to retrieve the token string from a cookie named "token".
See: https://github.com/go-chi/jwtauth/blob/1ff608193a049433794670a8c18fd739c5b2f236/jwtauth.go#L256
*/
func TokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie(AuthCookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}

/*
TokenFromHeader tries to retrieve the token string from the "Authorization"
request header: "Authorization: BEARER T".
See: https://github.com/go-chi/jwtauth/blob/1ff608193a049433794670a8c18fd739c5b2f236/jwtauth.go#L266
*/
func TokenFromHeader(r *http.Request) string {
	// Get token from authorization header.
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}

/*
TokenFromQuery tries to retrieve the token string from the "token" URI
query parameter.
See: https://github.com/go-chi/jwtauth/blob/1ff608193a049433794670a8c18fd739c5b2f236/jwtauth.go#L285
*/
func TokenFromQuery(r *http.Request) string {
	// Get token from query param named "jwt".
	return r.URL.Query().Get(AuthHttpParamName)
}

/*
This middleware is responsible for authenticating clients. A client can
provide a token from either a cookie (`TokenFromCookie()`), a header
(`TokenFromHeader()`), or from a URL parameter (`TokenFromQuery()`).
*/
func (amw authMiddleware) authMWHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Attempt to get the token, starting from the URL params
		token := TokenFromQuery(r)

		//If the token still isn't there, try the headers
		if token == "" {
			token = TokenFromHeader(r)
		}

		//If the token still isn't there, try the cookies
		if token == "" {
			token = TokenFromCookie(r)
		}

		//Still no token? Deny the request since there's no token. The middleware stops here
		if token == "" {
			httpu.HttpErrorAsJson(w, fmt.Errorf("auth; %s", ErrAuthNoTokenFound), http.StatusUnauthorized)
			return
		}

		//Get a byte array from the token and reject the token if the size is incorrect
		tbytes, derr := obj.Base64DecodeTok(token)
		if derr != nil || len(tbytes) != obj.TOKEN_SIZE_BYTES {
			httpu.HttpErrorAsJson(w, fmt.Errorf("auth; %s", ErrAuthBadTokenFormat), http.StatusUnauthorized)
			return
		}

		//Attempt to derive a token object from the input bytes
		tokObj := obj.TokenFromBytes(tbytes)
		if tokObj == nil {
			httpu.HttpErrorAsJson(w, fmt.Errorf("auth; %s", ErrAuthBadTokenFormat), http.StatusUnauthorized)
			return
		}

		/*
			Ensure the provided token isn't an invalid object for whatever reason. It's
			faster and safer to pre-check validity before anything else than to query for
			an invalid token. After this point, the token object itself is valid. Further
			checks can now occur before the database is queried for the subject's tokens.
		*/
		if !tokObj.Validate(true) {
			httpu.HttpErrorAsJson(w, fmt.Errorf("auth; %s", ErrAuthBadTokenFormat), http.StatusUnauthorized)
			return
		}

		//Check if the token has expired
		if tokExp := tokObj.GetExpiry(); tokObj.Expire && time.Now().After(tokExp) {
			httpu.HttpErrorAsJson(w, fmt.Errorf("auth; %s", ErrAuthExpiredToken), http.StatusUnauthorized)
			return
		}

		/*
			Ensure the token's scope is among those that are authorized. A token is considered
			valid for a route if the token's scope value is greater than or equal to the auth
			handler's lowest allowed scope. The scopes of an auth handler are pre-sorted in
			ascending order just after initialization.
		*/
		if !(tokObj.Scope >= amw.allowedScopes[0]) {
			httpu.HttpErrorAsJson(w, fmt.Errorf("auth; %s", ErrAuthUnauthorized), http.StatusUnauthorized)
			return
		}

		//Get the subject of the token
		tokSubject := tokObj.Subject

		//
		// -- BEGIN: Database Query
		//

		//Query the database for user tokens; if any errors occur, report them but do not alert the client of the specifics
		subjectTokens, dberr := crud.GetSTokens(amw.mclient, amw.rclient, r.Context(), tokSubject)
		if dberr != nil {
			httpu.HttpErrorAsJson(w, fmt.Errorf("auth; %s", ErrAuthGeneric), http.StatusUnauthorized)
			return
		}

		//If the query returns nothing, simply bail out; the subject has no tokens to begin with
		if len(subjectTokens) == 0 {
			httpu.HttpErrorAsJson(w, fmt.Errorf("auth; %s", ErrAuthNoTokenFound), http.StatusUnauthorized)
			return
		}

		//
		// -- END: Database Query
		//

		/*
			Check if the subject's token list includes the incoming token. If this check
			passes, the client is let through and the middleware finishes without error.
		*/
		if !slices.Contains(subjectTokens, token) {
			httpu.HttpErrorAsJson(w, fmt.Errorf("auth; %s", ErrAuthNoTokenFound), http.StatusUnauthorized)
			return
		}

		//Add headers to the request (auth subject and token scope)
		r.Header.Add(AuthHttpHeaderSubject, tokSubject.String())
		r.Header.Add(AuthHttpHeaderScope, strconv.Itoa(int(tokObj.Scope)))

		//Forward the request; authentication passed successfully
		next.ServeHTTP(w, r)
	})
}
