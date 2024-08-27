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
func GetRefreshTokFromCookie(w http.ResponseWriter, r *http.Request, failOnMissing bool) string {
	stok := util.StringFromCookie(r, token.RefreshTokenName)
	if stok == "" && failOnMissing {
		util.ErrResponse(
			AuthNoRToken,
			fmt.Errorf("refresh token is required"),
		).Respond(w)
		return ""
	}
	return stok
}

func IssueRefreshToken(user *user.User, ucoll *user.UserCollection, w http.ResponseWriter, persistent bool) {

}

/*
Attempts to decrypt and validate a refresh token. This is the first step
to run when attempting to reissue a refresh token.
*/
func DecodeRefreshToken() {

}
