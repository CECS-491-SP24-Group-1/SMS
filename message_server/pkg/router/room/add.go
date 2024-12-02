package room

import (
	"fmt"
	"net/http"

	"wraith.me/message_server/pkg/mw"
	"wraith.me/message_server/pkg/schema/user"
	"wraith.me/message_server/pkg/util"
)

// Handles incoming requests made to `PATCH /api/chat/room/{roomID}/add`.
func AddRoomRoute(w http.ResponseWriter, r *http.Request) {
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
		//Allow entry into the room
		room.AddMember(requestor.ID)

		//Save the chat room in the database
		_, err := rc.UpsertId(r.Context(), room.ID, room)
		if err != nil {
			util.ErrResponse(http.StatusInternalServerError, err).Respond(w)
			return
		}
	}

	fmt.Printf("user %s attempted to add room %s\n", requestor.ID, room.ID)

	//Respond with the created chat room using PayloadResponse
	util.PayloadOkResponse("Chat room created successfully", room).Respond(w)
}
