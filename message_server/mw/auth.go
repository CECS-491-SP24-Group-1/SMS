package mw

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"wraith.me/message_server/db"
	"wraith.me/message_server/db/mongoutil"
	"wraith.me/message_server/obj"
	"wraith.me/message_server/util/httpu"
)

var (
	//The name of the cookie to look for.
	AuthCookieName = "token"

	//The name of the HTTP query parameter to look for.
	AuthHttpParamName = "token"
)

// Holds the error messages.
var (
	ErrAuthUnauthorized   = errors.New("token is unauthorized")
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
func NewAuthMiddleware(allowedScopes []obj.TokenScope, mclient *mongo.Client, rclient *redis.Client) func(next http.Handler) http.Handler {
	//Get a struct object
	mw := authMiddleware{
		allowedScopes: allowedScopes,
		mclient:       mclient,
		rclient:       rclient,
	}

	//Return the instance
	slices.Sort(mw.allowedScopes)
	return mw.authMWHandler
}

/*
TokenFromCookie tries to retreive the token string from a cookie named "token".
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
		tbytes, derr := base64.StdEncoding.Strict().DecodeString(token)
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

		//Query the database and get the tokens of the subject; Redis is used here to cache the results
		var subjectTokens []string
		redisTokens, err := amw.rclient.LRange(r.Context(), tokSubject.String(), 0, -1).Result()
		if err == redis.Nil || len(redisTokens) == 0 {
			//Cache miss; query MongoDB for the tokens of the subject and add them to Redis for later
			//Define the aggregation to run
			query := bson.A{
				bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: tokSubject}}}},
				bson.D{
					{Key: "$project",
						Value: bson.D{
							{Key: "_id", Value: 0},
							{Key: "tokens", Value: 1},
						},
					},
				},
			}

			//Define the structure of the output; a list of string tokens
			type tokenRet struct {
				Tokens []string `bson:"tokens"`
			}

			//Execute the aggregation. Should return a singular `tokenRet` in the output array if the database has tokens for the subject
			ucoll := amw.mclient.Database(db.ROOT_DB).Collection(db.USERS_COLLECTION)
			var results []tokenRet
			if aerr := mongoutil.AggregateT(&results, ucoll, query, r.Context()); aerr != nil {
				//If there is a database aggregation error, deny the client by default for safety; do not report the error
				httpu.HttpErrorAsJson(w, fmt.Errorf("auth; %s", ErrAuthGeneric), http.StatusUnauthorized)
				fmt.Printf("[AUTH; MONGO AGGREGATION ERROR]: %s\n", aerr)
				return
			}

			//If the query returns nothing, simply bail out; the subject has no tokens to begin with
			if len(results) == 0 || len(results[0].Tokens) == 0 {
				httpu.HttpErrorAsJson(w, fmt.Errorf("auth; %s", ErrAuthNoTokenFound), http.StatusUnauthorized)
				return
			}

			//Copy the tokens to the subject tokens array
			subjectTokens = results[0].Tokens

			/*
				Add the subject ID and corresponding tokens array to Redis. Further
				requests with the same token subject will be honored by Redis instead
				of MongoDB (cache hit). If any errors occur, silently fail and do not
				alert the client. Authentication can still be done, but caching won't
				work.
			*/
			cacheSetResult := amw.rclient.RPush(r.Context(), tokSubject.String(), results[0].Tokens)
			if err := cacheSetResult.Err(); err != nil {
				fmt.Printf("[AUTH; REDIS CACHE RPUSH ERROR]: %s\n", err)
			}
		} else if err != nil {
			//If there is a cache retrieval error, deny the client by default for safety; do not report the error
			httpu.HttpErrorAsJson(w, fmt.Errorf("auth; %s", ErrAuthGeneric), http.StatusUnauthorized)
			fmt.Printf("[AUTH; REDIS ERROR]: %s\n", err)
			return
		} else {
			//Cache hit; copy the token list from Redis
			subjectTokens = redisTokens
		}

		/*
			Check if the subject's token list includes the incoming token. If this check
			passes, the client is let through and the middleware finishes without error.
		*/
		if !slices.Contains(subjectTokens, token) {
			httpu.HttpErrorAsJson(w, fmt.Errorf("auth; %s", ErrAuthNoTokenFound), http.StatusUnauthorized)
			return
		}

		fmt.Printf("TOKENS: %v\n", subjectTokens)
		fmt.Printf("Got token: %v\n", tokObj)

		//Forward the request; authentication passed successfully
		next.ServeHTTP(w, r)
	})
}
