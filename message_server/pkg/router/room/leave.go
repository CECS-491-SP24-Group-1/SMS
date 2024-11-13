package room

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/qiniu/qmgo"
	"wraith.me/message_server/pkg/mw"
	chatroom "wraith.me/message_server/pkg/schema/chat_room"
	"wraith.me/message_server/pkg/schema/user"
	"wraith.me/message_server/pkg/util"
)

// LeaveRoomRoute handles requests to `POST /api/chat/room/{roomID}/leave`.
func LeaveRoomRoute(w http.ResponseWriter, r *http.Request) {
	//Extract room ID from URL
	roomID := chi.URLParam(r, "roomID")
	rid, err := util.ParseUUIDv7(roomID)
	if err != nil {
		util.ErrResponse(http.StatusBadRequest, fmt.Errorf("bad room ID format; it must be a UUIDv7")).Respond(w)
		return
	}

	//Get the current user from the context
	requestor := r.Context().Value(mw.AuthCtxUserKey).(user.User)

	//Retrieve the room from the database
	var room chatroom.Room
	err = rc.FindID(r.Context(), rid).One(&room)
	if err != nil {
		code := http.StatusInternalServerError
		if qmgo.IsErrNoDocuments(err) {
			code = http.StatusNotFound
			err = fmt.Errorf("chat room with ID %s not found", rid)
		}
		util.ErrResponse(code, err).Respond(w)
		return
	}

	//Check if the user is a participant in the room
	if _, exists := room.Participants[requestor.ID]; !exists {
		util.ErrResponse(http.StatusForbidden, fmt.Errorf("you are not a member of this room")).Respond(w)
		return
	}

	//Remove the user from the participants list
	room.RemoveMember(requestor.ID)

	//Check if the room now has at least one member
	if room.Size() > 0 {
		//Upsert the room in the database
		_, err = rc.UpsertId(r.Context(), room.ID, room)
	} else {
		//Delete the room since nobody is left
		err = rc.RemoveId(r.Context(), room.ID)
		fmt.Printf("Room %s has no more members. Reaping...\n", roomID)
	}

	//Check if either of the operations failed
	if err != nil {
		util.ErrResponse(http.StatusInternalServerError,
			fmt.Errorf("failed to leave room with ID %s: %w", roomID, err),
		).Respond(w)
	}

	//Respond back with the ID of the room that was left
	util.OkResponse(fmt.Sprintf("successfully left the chat room with ID %s", roomID)).Respond(w)
}
