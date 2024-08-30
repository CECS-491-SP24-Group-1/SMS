package cauth

import (
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"wraith.me/message_server/config"
	"wraith.me/message_server/obj/ip_addr"
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
func AttemptRefreshAuth(w http.ResponseWriter, r *http.Request, env *config.Env,
	ucoll *user.UserCollection, failSilently bool) (usr *user.User, tid *util.UUID, err error) {
	//Rethrow errors into the HTTP response if any occur
	defer func() {
		if err != nil {
			if !failSilently {
				util.ErrResponse(http.StatusUnauthorized, err).Respond(w)
			} else {
				fmt.Printf(
					"error during refresh attempt for IP %s: %s\n",
					r.Host,
					err.Error(),
				)
			}
		}
	}()

	//Check if a user has a refresh token in the cookies
	rcookie := GetRefreshTokFromCookie(w, r, false)
	if rcookie == "" {
		err = fmt.Errorf("refresh token is required")
		return
	}

	//Decrypt the refresh token
	rtoken, err := token.Decrypt(rcookie, env.SK, env.ID, token.TokenTypeREFRESH)
	if err != nil {
		return
	}

	//Attempt to fetch a user from the database using the token subject field
	err = ucoll.Find(r.Context(), bson.M{"_id": rtoken.Subject}).One(&usr)
	if err != nil {
		return
	}

	//Ensure the user actually owns the token
	if !usr.HasTokenById(rtoken.ID.String()) {
		err = fmt.Errorf("none of the refresh tokens on-file match")
		return
	}

	//Nothing went wrong, so return normally
	//After this point, it is safe to assume that the user has proper auth
	tid = &rtoken.ID
	return
}

/*
Issues new access and refresh tokens for the user. This function ought to
run once the user has successfully authenticated either via a refresh token
or by solving a public key challenge.
*/
func PostAuth(
	w http.ResponseWriter, r *http.Request,
	usr *user.User, ucoll *user.UserCollection,
	cfg *token.TConfig, env *config.Env,
	persistent bool, tid *util.UUID,
) {
	//Delete the token from the user's list if one exists
	//The refresh token may also be "reused" at this step, but only a limited number of times
	if tid != nil {
		usr.RemoveToken(tid.String())
	}

	//Update the last IP and login fields of the user
	usr.LastIP = ip_addr.HttpIP2IPAddr(r.RemoteAddr)
	usr.LastLogin = util.NowMillis()

	//Issue an access and refresh token; this also updates the user in the database
	rtid, err := IssueRefreshToken(w, r, usr, ucoll, r.Context(), env, cfg, persistent)
	if err != nil {
		util.ErrResponse(http.StatusInternalServerError, err).Respond(w)
	}
	IssueAccessToken(w, r, usr, env, cfg, &rtid, persistent) //This should happen second; might want to tie access and refresh tokens together

	//Serialize the user's username and ID to a map
	payload := make(map[string]string)
	payload["id"] = usr.ID.String()
	payload["username"] = usr.Username

	//Respond back with the user's ID and username
	msg := "successfully logged in as"
	if tid != nil {
		msg = "successfully refreshed token for"
	}
	util.PayloadOkResponse(
		fmt.Sprintf("%s %s (id: %s)", msg, usr.Username, usr.ID.String()),
		payload,
	).Respond(w)
}
