package auth

import (
	"fmt"
	"net/http"

	"wraith.me/message_server/controller/csolver"
	"wraith.me/message_server/schema/user"
	"wraith.me/message_server/util"
)

/*
Handles incoming requests made to `POST /api/auth/login_req`. This is stage 1
of the login process.
*/
func RequestLoginUserRoute(w http.ResponseWriter, r *http.Request) {
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
	//Create a new stage 2 object plus database result
	loginVReq := csolver.LoginVerifyUser{}
	user := user.User{}

	//Run pre-flight checks
	if !csolver.PreFlight(&loginVReq, &user, uc, w, r) {
		return
	}

	//Verify the public key challenge
	//After this point, it is safe to assume that a user is authorized to login
	tok, err := csolver.VerifyPKChallenge(loginVReq, env)
	if err != nil {
		util.ErrResponse(http.StatusForbidden, err)
		return
	}

	//TODO: mark user as verified and issue a login token here

	fmt.Printf("verif_pk: %+v\n", tok)
	util.PayloadOkResponse("", "ok").Respond(w)
}
