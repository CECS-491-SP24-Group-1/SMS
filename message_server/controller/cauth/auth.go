package cauth

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"wraith.me/message_server/config"
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
func AttemptRefresh(w http.ResponseWriter, r *http.Request, env *config.Env,
	ucoll *user.UserCollection, failSilently bool) (usr *user.User, err error,
) {
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

	//Nothing went wrong, so return normally
	return
}

/*
Issues new access and refresh tokens for the user. This function ought to
run once the user has successfully authenticated either via a refresh token
or by solving a public key challenge.
*/
func PostAuth(
	w http.ResponseWriter, ctx context.Context,
	usr *user.User, ucoll *user.UserCollection,
	cfg *token.TConfig, env *config.Env,
	persistent bool, newToken bool,
) {
	//TODO: delete existing token here
	if !newToken {
	}

	//Issue an access and refresh token; this also updates the user in the database
	IssueAccessToken(w, usr, env, cfg, persistent)
	err := IssueRefreshToken(w, usr, ucoll, ctx, env, cfg, persistent)
	if err != nil {
		util.ErrResponse(http.StatusInternalServerError, err).Respond(w)
	}

	//Serialize the user's username and ID to a map
	payload := make(map[string]string)
	payload["id"] = usr.ID.String()
	payload["username"] = usr.Username

	//Respond back with the user's ID and username
	msg := "successfully logged in as"
	if !newToken {
		fmt.Println("dshahahda")
		msg = "successfully refreshed token for"
	}
	util.PayloadOkResponse(
		fmt.Sprintf("%s %s <id: %s>", msg, usr.Username, usr.ID.String()),
		payload,
	).Respond(w)
}
