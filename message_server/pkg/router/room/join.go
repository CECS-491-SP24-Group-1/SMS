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

// Handles incoming requests made to `GET /api/chat/room/{roomID}/join`.
func JoinRoomRoute(w http.ResponseWriter, r *http.Request) {
	//Get the room from the request params
	room := getRoomFromQuery(w, r)
	if room == nil {
		return
	}

	//Get the requestor's info
	requestor := r.Context().Value(mw.AuthCtxUserKey).(user.User)

	//Ensure the current user is allowed to join the room
	//TODO: handle perms for this later; let them anyway
	if !room.HasMember(requestor.ID) {
		//Allow entry into the room
		room.AddMember(requestor.ID)

		//Save the chat room in the database
		_, err := rc.UpsertId(r.Context(), room.ID, room)
		if err != nil {
			util.ErrResponse(http.StatusInternalServerError, err).Respond(w)
			return
		}
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
