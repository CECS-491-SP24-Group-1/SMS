package room

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/qiniu/qmgo"
	chatroom "wraith.me/message_server/pkg/schema/chat_room"
	"wraith.me/message_server/pkg/util"
)

// Derives a valid `Room` object from the request parameters.
func getRoomFromQuery(w http.ResponseWriter, r *http.Request) *chatroom.Room {
	//Get the ID of the chat room
	roomID := chi.URLParam(r, "roomID")
	rid, err := util.ParseUUIDv7(roomID)
	if err != nil {
		util.ErrResponse(
			http.StatusBadRequest,
			fmt.Errorf("bad room ID format; it must be a UUIDv7"),
		).Respond(w)
		fmt.Printf("bad room ID\n")
		return nil
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
		return nil
	}

	//Return the room object
	return &room
}
