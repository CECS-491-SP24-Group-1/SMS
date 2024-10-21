package wschat

import (
	"fmt"

	"github.com/olahol/melody"
	"wraith.me/message_server/pkg/util"
)

// Handles new connections to the chat server.
func (w *Server) handleConnect(s *melody.Session) {
	//Get the room's ID
	roomID, ok := s.Request.Context().Value(WSChatCtxRoomIDKey).(util.UUID)
	if !ok {
		s.Write([]byte("Invalid room ID format"))
		s.Close()
		return
	}

	//Get an existing room or create a new one and add the new user
	room := w.getOrCreateRoom(roomID)
	room.AddSession(s)

	//Send the MOTD to the connecting user and announce the membership change
	s.Write([]byte("Welcome to the room!"))
	announceMembershipChange(s, room, true)
}

// Handles disconnections from the chat server.
func (w *Server) handleDisconnect(s *melody.Session) {
	//Get the room's ID
	roomUUID, ok := s.Request.Context().Value(WSChatCtxRoomIDKey).(util.UUID)
	if !ok {
		return
	}

	//Get the room instance and eject the session if it's non-null
	room := w.getRoom(roomUUID)
	if room != nil {
		//Remove the current session handler for the user
		room.RemoveSession(s)

		//Broadcast the membership change event
		announceMembershipChange(s, room, false)

		//Eject the room from the handler if the last person left
		if room.IsEmpty() {
			w.removeRoom(roomUUID)
		}
	}
}

// Announces the membership info to the room.
func announceMembershipChange(s *melody.Session, room *Room, joined bool) {
	//Construct the initial message
	initial := fmt.Sprintf(
		"A user %s the room. There are %d members online.",
		util.If(joined, "joined", "left"),
		room.Size(),
	)

	//Broadcast to the room
	room.Broadcast([]byte(initial), s)
}
