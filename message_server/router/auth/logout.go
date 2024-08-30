package auth

import (
	"fmt"
	"net/http"

	"wraith.me/message_server/controller/cauth"
	"wraith.me/message_server/obj/token"
	"wraith.me/message_server/util"
)

// Handles incoming requests made to `POST /api/auth/logout`.
func LogoutRoute(w http.ResponseWriter, r *http.Request) {
	//Attempt to authenticate via the user's refresh token
	user, tid, err := cauth.AttemptRefreshAuth(w, r, env, uc, false)

	//Run post-auth if the process succeeded
	//The refresh attempt will auto-respond if something goes wrong
	if user != nil && err == nil {
		//Delete the token from the list of the user's tokens
		user.RemoveToken(tid.String())

		//Upsert the corresponding document in the database
		_, err := uc.UpsertId(r.Context(), user.ID, user)
		if err != nil {
			util.ErrResponse(http.StatusInternalServerError, err).Respond(w)
		} else {
			/*
				"Delete" the tokens by replacing them with empty & expired copies.
				The web browser will vacuum up these cookies automatically, but only
				if the domain and path match exactly. Go has no in-built function to
				do this, so each cookie is replaced instead. Research into Express's
				method reveals a similar approach, so this way is most likely the
				correct way to handle deletion of cookies.
			*/
			util.DeleteCookie(w, token.AccessTokenName, cfg.Token.Domain, cfg.Token.AccessCookiePath)
			util.DeleteCookie(w, token.AccessTokenExprName, cfg.Token.Domain, cauth.ExprCookiePath)
			util.DeleteCookie(w, token.RefreshTokenName, cfg.Token.Domain, cfg.Token.RefreshCookiePath)
			util.DeleteCookie(w, token.RefreshTokenExprName, cfg.Token.Domain, cauth.ExprCookiePath)

			//Respond back that the logout was successful
			util.OkResponse(
				fmt.Sprintf(
					"%s (id: %s) logged out successfully; have a secure day!",
					user.Username, user.ID.String(),
				)).Respond(w)
		}
	}
}
