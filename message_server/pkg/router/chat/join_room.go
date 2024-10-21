package chat

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"wraith.me/message_server/pkg/util"
	"wraith.me/message_server/pkg/ws/wschat"
)

// Handles incoming requests made to `GET /api/chat/room/{roomID}`.
func JoinRoomRoute(w http.ResponseWriter, r *http.Request) {
	log.Println("Request URL:", r.URL.String())

	//Get the ID of the chat room
	roomID := chi.URLParam(r, "roomID")
	rid, err := util.ParseUUIDv7(roomID)
	if err != nil {
		util.ErrResponse(
			http.StatusBadRequest,
			fmt.Errorf("bad room ID format; it must be a UUIDv7"),
		).Respond(w)
		return
	}

	//TODO: check if a user is a member of the room and if the room even exists

	//Set the room ID in the request context and handle the connection
	r = r.WithContext(context.WithValue(r.Context(), wschat.WSChatCtxRoomIDKey, rid))
	mel.GetMelody().HandleRequest(w, r)
}
