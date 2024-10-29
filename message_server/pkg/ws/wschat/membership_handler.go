package wschat

import (
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
	room := w.getOrCreateRoom(*roomID, getParticipants(s))

	//Reject the session if the user is already in the room
	if room.HasUser(uinfo.ID) {
		s.Write([]byte("You are already in the room"))
		s.Close()
		return
	}

	//Get the number of users currently in the room
	currentUsers := room.Size()

	//Add the user to the room
	room.AddSession(s, uinfo)

	//Announce the membership change
	announceMembershipChange(s, room, currentUsers)
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
		//Get the number of users currently in the room
		currentUsers := room.Size()

		//Remove the current session handler for the user
		room.RemoveSession(s)

		//Broadcast the membership change event
		announceMembershipChange(s, room, currentUsers)

		//Eject the room from the handler if the last person left
		if room.IsEmpty() {
			w.removeRoom(*roomUUID)
		}
	}
}

// Announces the membership info to the room.
func announceMembershipChange(_ *melody.Session, room *WSRoom, oldSize int) {
	//Get the numbers of members in the room after the change
	newSize := room.Size()

	//Get the type of membership change event
	typ := util.If(newSize > oldSize, chat.TypeJOINEVENT, chat.TypeQUITEVENT)

	//Construct the inner membership change event message
	content := chat.MembershipChange{
		Old: oldSize,
		New: newSize,
	}

	//Create the message
	msg := chat.NewMessageTyp(string(content.JSON()), room.ID, room.ID, typ)

	//Broadcast to the room
	room.Broadcast(msg.JSON())
}
