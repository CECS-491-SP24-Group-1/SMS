package chat

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/olahol/melody"
)

// TODO: tmp

/*
Handles incoming requests made to `GET /api/chat/room/{roomID}`.
*/
func ChatRoomRoute(w http.ResponseWriter, r *http.Request) {
	log.Println("Request URL:", r.URL.String())

	//Get the ID of the chat room
	roomID := chi.URLParam(r, "roomID")
	log.Printf("url param: %s\n", roomID)

	if roomID == "" {
		http.Error(w, "roomID parameter is missing", http.StatusBadRequest)
		return
	}

	//Get the Melody instance
	m := mel.GetMelody()

	//Add the handler
	//TODO: move this out eventually
	m.HandleMessage(func(s *melody.Session, msg []byte) {
		log.Printf("<%s::%s>: %s\n", r.RemoteAddr, roomID, string(msg))
		m.Broadcast(msg)
	})

	m.HandleRequest(w, r)
}
