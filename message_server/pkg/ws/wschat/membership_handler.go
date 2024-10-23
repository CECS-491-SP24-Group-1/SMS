package wschat

import (
	"fmt"

	"github.com/olahol/melody"
	"wraith.me/message_server/pkg/http_types/ws/chat"
	"wraith.me/message_server/pkg/util"
)

// Handles new connections to the chat server.
func (w *Server) handleConnect(s *melody.Session) {
	//Get the room's ID
	roomID := getRoomID(s)
	if roomID == nil {
		//Emit an error message and eject the client
		s.Write([]byte("Invalid room ID format"))
		s.Close()
		return
	}

	//Get the ID of the user who is trying to connect
	userID := getUserID(s)
	if userID == nil {
		return
	}

	//Create the user data object for the room
	uinfo := &UserData{
		ID: *userID,
	}

	//Get an existing room or create a new one
	room := w.getOrCreateRoom(*roomID)

	//Reject the session if the user is already in the room
	if room.HasUser(uinfo.ID) {
		s.Write([]byte("You are already in the room"))
		s.Close()
		return
	}

	//Add the user to the room
	room.AddSession(s, uinfo)

	//Send the MOTD to the connecting user and announce the membership change
	s.Write([]byte("Welcome to the room!"))
	announceMembershipChange(s, room, true)
}

// Handles disconnections from the chat server.
func (w *Server) handleDisconnect(s *melody.Session) {
	//Get the room's ID
	roomUUID := getRoomID(s)
	if roomUUID == nil {
		return
	}

	//Get the room instance and eject the session if it's non-null
	room := w.getRoom(*roomUUID)
	if room != nil {
		//Remove the current session handler for the user
		room.RemoveSession(s)

		//Broadcast the membership change event
		announceMembershipChange(s, room, false)

		//Eject the room from the handler if the last person left
		if room.IsEmpty() {
			w.removeRoom(*roomUUID)
		}
	}
}

// Announces the membership info to the room.
func announceMembershipChange(s *melody.Session, room *WSRoom, joined bool) {
	//Construct the initial message
	content := fmt.Sprintf(
		"A user %s the room. There are %d members online.",
		util.If(joined, "joined", "left"),
		room.Size(),
	)
	msg := chat.NewMessage(content, room.ID, room.ID)

	//Broadcast to the room
	room.Broadcast(msg.JSON(), s)
}
