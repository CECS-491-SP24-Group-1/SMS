package user

import (
	"net/http"

	"wraith.me/message_server/pkg/http_types/response"
	"wraith.me/message_server/pkg/mw"
	"wraith.me/message_server/pkg/schema/user"
	"wraith.me/message_server/pkg/util"
)

// Handles incoming requests made to `GET /api/user/me` and `GET /api/user/`.
func HandleMyInfoRoute(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(mw.AuthCtxUserKey).(user.User)
	sendUserInfo(w, user)
}

// Sends user info to an HTTP response; truncated.
func sendUserInfo(w http.ResponseWriter, user user.User) {
	info := response.FromUser(user)
	util.PayloadOkResponse("", info).Respond(w)
}
