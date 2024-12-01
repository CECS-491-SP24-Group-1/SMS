package wschat

import "github.com/olahol/melody"

// Handles messages sent to the chat server by participants.
func (w *Server) handleMessage(s *melody.Session, msg []byte) {
	//Get the ID of the room
	roomUUID := getRoomID(s)
	if roomUUID == nil {
		return
	}

	room := w.GetRoom(*roomUUID)
	if room != nil && room.HasSession(s) {
		if isValidMessage(msg) {
			room.Broadcast(msg, nil)
		} else {
			s.Write([]byte("Invalid message"))
		}
	}
}

func isValidMessage(msg []byte) bool {
	return len(msg) > 0
}
