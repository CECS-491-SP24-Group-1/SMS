package auth

import (
	"net/http"

	"wraith.me/message_server/controller/cauth"
)

// Handles incoming requests made to `POST /api/auth/refresh`.
func RefreshTokenRoute(w http.ResponseWriter, r *http.Request) {
	//Attempt to authenticate via the user's refresh token
	user, tid, err := cauth.AttemptRefreshAuth(w, r, env, uc, false)

	//Run post-auth if the process succeeded
	//The refresh attempt will auto-respond if something goes wrong
	if user != nil && err == nil {
		cauth.PostAuth(w, r.Context(), user, uc, &cfg.Token, env, true, tid)
		return
	}
}
