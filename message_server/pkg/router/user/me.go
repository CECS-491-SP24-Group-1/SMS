package user

import (
	"net/http"
	"fmt"

	"wraith.me/message_server/pkg/mw"
	"wraith.me/message_server/pkg/schema/user"
	"wraith.me/message_server/pkg/util"

)

func getMyInfo(w http.ResponseWriter, r *http.Request) {

	// Get the requestor's info
	user := r.Context().Value(mw.AuthCtxUserKey).(user.User)

	util.PayloadOkResponse(fmt.Sprintf("user retrieved successfully"), user).Respond(w)

}