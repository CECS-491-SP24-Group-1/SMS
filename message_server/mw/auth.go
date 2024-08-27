package mw

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"wraith.me/message_server/config"
	"wraith.me/message_server/obj"
	"wraith.me/message_server/obj/token"
	cr "wraith.me/message_server/redis"
	"wraith.me/message_server/schema/user"
	"wraith.me/message_server/util"
)

var (
	//The name of the cookie to look for.
	AuthCookieName = token.AccessTokenName

	//The name of the HTTP query parameter to look for.
	AuthHttpParamName = token.AccessTokenName

	//The name of the auth subject header to send.
	AuthHttpHeaderSubject = "X-Auth-For"

	//The key of the user object that's passed via `r.Context`.
	AuthCtxUserKey = obj.CtxKey{S: "ReqUser"}
)

// Holds the error messages.
var (
	ErrAuthUnauthorized = errors.New("token is not authorized for this route")
	ErrAuthNoTokenFound = errors.New("no token found")
	ErrAuthGeneric      = errors.New("authentication error")
	ErrAuthNotFound     = errors.New("user not found with ID %s")
)

type authMiddleware struct {
	//allowedScopes []token.TokenScope //The scopes for which the token is valid, sorted in increasing order.
	//mclient       *mongo.Client      //The MongoDB database client.
	rclient *redis.Client //The Redis database client.

	//The user collection to use.
	ucoll *user.UserCollection

	//The secrets of the server including ID and encryption key.
	secrets *config.Env
}

// Returns a new handler for the authentication middleware.
func NewAuthMiddleware(secrets *config.Env) func(next http.Handler) http.Handler {
	//Get a struct object
	mw := authMiddleware{
		//allowedScopes: allowedScopes,
		//mclient:       db.GetInstance().GetClient(),
		rclient: cr.GetInstance().GetClient(),
		ucoll:   user.GetCollection(),
		secrets: secrets,
	}

	//Return the instance
	//slices.Sort(mw.allowedScopes)
	return mw.authMWHandler
}

/*
TokenFromHeader tries to retrieve the token string from the "Authorization"
request header: "Authorization: BEARER T".
See: https://github.com/go-chi/jwtauth/blob/1ff608193a049433794670a8c18fd739c5b2f236/jwtauth.go#L266
*/
func TokenFromHeader(r *http.Request) string {
	//Get token from authorization header
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}

/*
This middleware is responsible for authenticating clients. A client can
provide a token from either a cookie (`TokenFromCookie()`), a header
(`TokenFromHeader()`), or from a URL parameter (`TokenFromQuery()`).
*/
func (amw authMiddleware) authMWHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Attempt to get the token, starting from the URL params
		tok := util.StringFromQuery(r, AuthHttpParamName)

		//If the token still isn't there, try the headers
		if tok == "" {
			tok = TokenFromHeader(r)
		}

		//If the token still isn't there, try the cookies
		if tok == "" {
			tok = util.StringFromCookie(r, AuthCookieName)
		}

		//Still no token? Deny the request since there's no token. The middleware stops here
		if tok == "" {
			util.ErrResponse(
				http.StatusUnauthorized,
				fmt.Errorf("auth; %s", ErrAuthNoTokenFound),
			).Respond(w)
			return
		}

		//Decrypt and validate the authentication token
		tokObj, err := token.Decrypt(
			tok,
			amw.secrets.SK,
			amw.secrets.ID,
			token.TokenTypeACCESS,
		)
		if err != nil {
			util.ErrResponse(
				http.StatusUnauthorized,
				fmt.Errorf("auth; %s", err),
			).Respond(w)
			return
		}

		/*
				Ensure the token's scope is among those that are authorized. A token is considered
				valid for a route if the token's scope value is greater than or equal to the auth
				handler's lowest allowed scope. The scopes of an auth handler are pre-sorted in
				ascending order just after initialization.
			* /
			if !(tokObj.Scope >= amw.allowedScopes[0]) {
				util.ErrResponse(
					http.StatusUnauthorized,
					fmt.Errorf("auth; %s", ErrAuthUnauthorized),
				).Respond(w)
				return
			}
		*/

		//Get the subject of the token
		tokSubject := tokObj.Subject

		//
		// -- BEGIN: Database Query
		//

		//Prepare an object to hold the database result
		var user user.User

		//Query the database for the token's subject
		//After this point, assuming nothing goes wrong, the user is considered authorized to continue
		err = amw.ucoll.Find(
			r.Context(),
			bson.M{"_id": tokSubject},
		).One(&user)

		//Check if something went wrong during the query
		if err != nil {
			//Check if the error has to do with a lack of documents
			code := http.StatusUnauthorized
			desc := ErrAuthGeneric
			if errors.Is(err, mongo.ErrNoDocuments) {
				//Change the error to be a 404
				code = http.StatusNotFound
				desc = fmt.Errorf(ErrAuthNotFound.Error(), tokSubject.String())
			} else {
				fmt.Printf("auth err for client %s; %s\n", r.RemoteAddr, err)
			}

			//Respond back with the error
			util.ErrResponse(
				code,
				fmt.Errorf("auth; %s", desc),
			).Respond(w)
			return
		}
		//query := bson.D{{Key: "tokens", Value: bson.D{{Key: "$in", Value: bson.A{tok}}}}}

		//
		// -- END: Database Query
		//

		//Add headers to the request (auth subject and token scope)
		r.Header.Add(AuthHttpHeaderSubject, tokSubject.String())

		//Add the user to the request context
		//https://go.dev/blog/context#TOC_3.2.
		ctx := context.WithValue(r.Context(), AuthCtxUserKey, user)
		r = r.WithContext(ctx)

		//Forward the request; authentication passed successfully
		next.ServeHTTP(w, r)
	})
}
