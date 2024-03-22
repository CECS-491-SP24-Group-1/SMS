package mw

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
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
	ErrAuthExpiredToken   = errors.New("token is expired")
	ErrAuthNoTokenFound   = errors.New("no token found")
	ErrAuthBadTokenFormat = errors.New("token format is incorrect")
	ErrAuthGeneric        = errors.New("authentication error")
)

type authMiddleware struct {
	allowedScopes []obj.TokenScope
	mclient       *mongo.Client
	rclient       *redis.Client
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

		//Get the subject of the token
		tokSubject := tokObj.Subject

		//Ensure the provided token isn't expired or invalid
		//It's faster to pre-check validity than to query an invalid token
		if tokExp := tokObj.GetExpiry(); tokObj.Expire && time.Now().After(tokExp) {
			httpu.HttpErrorAsJson(w, fmt.Errorf("auth; %s", ErrAuthExpiredToken), http.StatusUnauthorized)
			return
		}

		//Ensure the token object is valid
		if !tokObj.Validate() {
			httpu.HttpErrorAsJson(w, fmt.Errorf("auth; %s", ErrAuthBadTokenFormat), http.StatusUnauthorized)
			return
		}

		//Query the database and get the tokens of the subject; Redis is used here to cache the results
		subjectTokens := make([]string, 0)
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
			fmt.Println("[auth] Cache Miss!")

			//Add the subject ID and corresponding tokens array to Redis.
			//Further requests with the same token subject will be honored by Redis instead.
			//If any errors occur, silently fail. Authentication can still be done, but caching won't work.
			cacheSetResult := amw.rclient.RPush(r.Context(), tokSubject.String(), results[0].Tokens)
			if err := cacheSetResult.Err(); err != nil {
				fmt.Printf("[AUTH; REDIS CACHE RPUSH ERROR]: %s\n", err)
			}
		} else if err != nil {
			//If there is a cache error, deny the client by default for safety; do not report the error
			httpu.HttpErrorAsJson(w, fmt.Errorf("auth; %s", ErrAuthGeneric), http.StatusUnauthorized)
			fmt.Printf("[AUTH; REDIS ERROR]: %s\n", err)
		} else {
			//Cache hit; copy the token list from Redis
			fmt.Println("[auth] Cache Hit!")
			subjectTokens = redisTokens
		}

		fmt.Printf("TOKENS: %v\n", subjectTokens)

		//Check if the token is expired
		//TODO: get the token from the user db first
		//fmt.Printf("Tok expiry: %s\n", tokObj.GetExpiry())
		//fmt.Printf("Tok subject: %s\n", tokSubject)

		if tokExp := tokObj.GetExpiry(); time.Now().After(tokExp) {
			httpu.HttpErrorAsJson(w, fmt.Errorf("auth; %s", ErrAuthExpiredToken), http.StatusUnauthorized)
			return
		}

		fmt.Printf("Got token: %v\n", tokObj)

		//Forward the request; authentication passed successfully
		next.ServeHTTP(w, r)
	})
}
