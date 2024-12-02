package room

import (
	"context"
	"fmt"
	"net/http"

	"wraith.me/message_server/pkg/mw"
	"wraith.me/message_server/pkg/schema/user"
	"wraith.me/message_server/pkg/util"
	"wraith.me/message_server/pkg/ws/wschat"
)

// Handles incoming requests made to `GET /api/chat/room/{roomID}`.
func JoinRoomRoute(w http.ResponseWriter, r *http.Request) {
	//Get the room from the request params
	room := getRoomFromQuery(w, r)
	if room == nil {
		return
	}

	//Get the requestor's info
	requestor := r.Context().Value(mw.AuthCtxUserKey).(user.User)

	//Ensure the current user is allowed to join the room
	//TODO: handle perms for this later; let them anywayy
	if !room.HasMember(requestor.ID) {
		util.ErrResponse(
			http.StatusForbidden,
			fmt.Errorf("you are not a member of this room"),
		).Respond(w)
		fmt.Printf("user %s denied entry to room %s; not a member\n", requestor.ID, room.ID)
		return
	}

	//Create the context object
	ctx := wschat.Context{
		RoomID:       room.ID,
		MemberID:     requestor.ID,
		Participants: room.Participants,
	}

	fmt.Printf("user %s attempted to join room %s\n", requestor.ID, room.ID)

	//Set the request context and handle the connection
	r = r.WithContext(context.WithValue(r.Context(), wschat.WSChatCtxObjKey, ctx))
	mel.GetMelody().HandleRequest(w, r)
}
