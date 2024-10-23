package room

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/qiniu/qmgo"
	"wraith.me/message_server/pkg/mw"
	chatroom "wraith.me/message_server/pkg/schema/chat_room"
	"wraith.me/message_server/pkg/schema/user"
	"wraith.me/message_server/pkg/util"
	"wraith.me/message_server/pkg/ws/wschat"
)

// Handles incoming requests made to `GET /api/chat/room/{roomID}`.
func JoinRoomRoute(w http.ResponseWriter, r *http.Request) {
	//Get the ID of the chat room
	roomID := chi.URLParam(r, "roomID")
	rid, err := util.ParseUUIDv7(roomID)
	if err != nil {
		util.ErrResponse(
			http.StatusBadRequest,
			fmt.Errorf("bad room ID format; it must be a UUIDv7"),
		).Respond(w)
		fmt.Printf("bad room ID\n")
		return
	}

	//Get the room info from the database
	var room chatroom.Room
	err = rc.FindID(r.Context(), rid).One(&room)
	if err != nil {
		//Handle 404s differently
		code := http.StatusInternalServerError
		if qmgo.IsErrNoDocuments(err) {
			code = http.StatusNotFound
			err = fmt.Errorf("cannot find chat room with ID %s", rid)
			fmt.Printf("%s is not a valid room\n", rid)
		}
		util.ErrResponse(code, err).Respond(w)
		return
	}

	//Ensure the current user is allowed to join the room
	requestor := r.Context().Value(mw.AuthCtxUserKey).(user.User)
	if _, ok := room.Participants[requestor.ID]; !ok {
		util.ErrResponse(
			http.StatusForbidden,
			fmt.Errorf("you are not a member of this room"),
		).Respond(w)
		fmt.Printf("user %s denied entry to room %s; not a member\n", requestor.ID, rid)
		return
	}

	//Create the context object
	ctx := wschat.Context{
		RoomID:   rid,
		MemberID: requestor.ID,
	}

	fmt.Printf("user %s attempted to join room %s\n", requestor.ID, rid)

	//Set the room ID and user ID in the request context and handle the connection
	r = r.WithContext(context.WithValue(r.Context(), wschat.WSChatCtxObjKey, ctx))
	mel.GetMelody().HandleRequest(w, r)
}
