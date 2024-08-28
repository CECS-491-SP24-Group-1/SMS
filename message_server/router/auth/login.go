package auth

import (
	"context"
	"fmt"
	"net/http"

	"wraith.me/message_server/controller/cauth"
	"wraith.me/message_server/controller/csolver"
	"wraith.me/message_server/schema/user"
	"wraith.me/message_server/util"
)

//TODO: attempt to read and verify a refresh token in the first step; verification is deemed unnecessary if this is the case

/*
Handles incoming requests made to `POST /api/auth/login_req`. This is stage 1
of the login process.
*/
func RequestLoginUserRoute(w http.ResponseWriter, r *http.Request) {
	//Skip straight to the post-login process if the user possesses a refresh token
	if user, err := cauth.AttemptRefresh(w, r, true); user != nil && err == nil {
		PostLogin(w, r.Context(), user, true, false)
		return
	}

	//Create a new stage 1 object plus database result
	loginReq := csolver.LoginUser{}
	user := user.User{}

	//Run pre-flight checks
	if !csolver.PreFlight(&loginReq, &user, uc, w, r) {
		return
	}

	//Check the user's flags to ensure they can actually sign-in
	//Their email must be verified before this may occur
	if !user.Flags.EmailVerified {
		util.ErrResponse(
			http.StatusForbidden,
			fmt.Errorf("unverified email"),
		).Respond(w)
		return
	}

	//Create a public key challenge using the user's info
	loginTok := csolver.IssuePKChallenge(user, env)

	//Send the token to the user
	util.PayloadOkResponse(
		"",
		loginTok,
	).Respond(w)
}

/*
Handles incoming requests made to `POST /api/auth/login_verify`. This is stage
2 of the login process.
*/
func VerifyLoginUserRoute(w http.ResponseWriter, r *http.Request) {
	//Skip straight to the post-login process if the user possesses a refresh token
	if user, err := cauth.AttemptRefresh(w, r, true); user != nil && err == nil {
		PostLogin(w, r.Context(), user, true, false)
		return
	}

	//Create a new stage 2 object plus database result
	loginVReq := csolver.LoginVerifyUser{}
	user := user.User{}

	//Run pre-flight checks
	if !csolver.PreFlight(&loginVReq, &user, uc, w, r) {
		return
	}

	//Verify the public key challenge
	//After this point, it is safe to assume that a user is authorized to login
	_, err := csolver.VerifyPKChallenge(loginVReq, env)
	if err != nil {
		util.ErrResponse(http.StatusForbidden, err)
		return
	}

	//Mark the user as PK verified and run post-login stuff
	user.MarkPKVerified()
	PostLogin(w, r.Context(), &user, true, true)
}

// Runs the logic after a successful login verification.
func PostLogin(w http.ResponseWriter, ctx context.Context, usr *user.User, persistent bool, newToken bool) {
	//Issue an access and refresh token; this also updates the user in the database
	cauth.IssueAccessToken(w, usr, env, &cfg.Token, persistent)
	err := cauth.IssueRefreshToken(w, usr, uc, ctx, env, &cfg.Token, persistent)
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
		msg = "successfully refreshed token for"
	}
	util.PayloadOkResponse(
		fmt.Sprintf("%s %s <id: %s>", msg, usr.Username, usr.ID.String()),
		payload,
	).Respond(w)
}
