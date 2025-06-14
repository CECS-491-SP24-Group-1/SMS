package auth

import (
	"fmt"
	"net/http"

	"wraith.me/message_server/pkg/controller/cauth"
	"wraith.me/message_server/pkg/controller/csolver"
	"wraith.me/message_server/pkg/http_types/response"
	"wraith.me/message_server/pkg/schema/user"
	"wraith.me/message_server/pkg/util"
)

/*
Handles incoming requests made to `POST /api/auth/login_req`. This is stage 1
of the login process.
*/
func RequestLoginUserRoute(w http.ResponseWriter, r *http.Request) {
	//Skip straight to the post-login process if the user possesses a refresh token
	if user, tid, err :=
		cauth.AttemptRefreshAuth(w, r, env, uc, true); user != nil && err == nil {
		fmt.Println("post auth in stage1")
		cauth.PostAuth(w, r, user, uc, &cfg.Token, env, true, tid)
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
		response.LoginReq{Token: loginTok},
	).Respond(w)
}

/*
Handles incoming requests made to `POST /api/auth/login_verify`. This is stage
2 of the login process.
*/
func VerifyLoginUserRoute(w http.ResponseWriter, r *http.Request) {
	//Skip straight to the post-login process if the user possesses a refresh token
	if user, tid, err :=
		cauth.AttemptRefreshAuth(w, r, env, uc, true); user != nil && err == nil {
		fmt.Println("post auth in stage2")
		cauth.PostAuth(w, r, user, uc, &cfg.Token, env, true, tid)
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
	_, err := csolver.VerifyPKChallenge(loginVReq, env, r)
	if err != nil {
		util.ErrResponse(http.StatusForbidden, err).Respond(w)
		return
	}

	//Mark the user as PK verified and run post-login stuff
	user.MarkPKVerified()
	cauth.PostAuth(w, r, &user, uc, &cfg.Token, env, true, nil)
}
