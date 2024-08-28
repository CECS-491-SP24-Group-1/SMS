package cauth

import (
	"fmt"
	"net/http"

	"wraith.me/message_server/obj/token"
	"wraith.me/message_server/schema/user"
	"wraith.me/message_server/util"
)

var (
	//The HTTP status code to emit when the refresh token wasn't found.
	AuthNoRToken = http.StatusForbidden
)

/*
Attempts to get the refresh token from the request cookies. This is the
first thing to do when attempting to issue a refresh token, as sessions
should be reused whenever possible.
*/
func GetRefreshTokFromCookie(w http.ResponseWriter, r *http.Request, failSilently bool) string {
	stok := util.StringFromCookie(r, token.RefreshTokenName)
	return stok
}

/*
Attempts to authenticate a user via a refresh token in the request cookies.
This function is to run before any login request/validation as well as during
a token refresh operation. In the case of logins, this function allows for
quicker turnarounds if a user has a valid refresh token, as it allows the user
to bypass the whole login process. If any errors occur during the login process,
they should occur silently and without disturbance to the process. Token
refreshes, on the other hand, should not fail silently.
*/
func AttemptRefresh(w http.ResponseWriter, r *http.Request, failSilently bool) (usr *user.User, err error) {
	//Rethrow errors into the HTTP response if any occur
	defer func() {
		if !failSilently && err != nil {
			util.ErrResponse(http.StatusUnauthorized, err).Respond(w)
		}
	}()

	//Check if a user has a refresh token in the cookies
	rtoken := GetRefreshTokFromCookie(w, r, false)
	if rtoken == "" {
		err = fmt.Errorf("refresh token is required")
		return
	}

	fmt.Printf("has refresh token %s\n", rtoken)
	return nil, nil
}

/*
Attempts to decrypt and validate a refresh token. This is the first step
to run when attempting to reissue a refresh token.
*/
func DecodeRefreshToken() {

}
