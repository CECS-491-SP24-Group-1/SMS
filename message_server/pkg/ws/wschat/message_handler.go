package wschat

import (
	"github.com/olahol/melody"
	"wraith.me/message_server/pkg/util"
)

// Handles messages sent to the chat server by participants.
func (w *Server) handleMessage(s *melody.Session, msg []byte) {
	roomUUID, ok := s.Request.Context().Value(WSChatCtxRoomIDKey).(util.UUID)
	if !ok {
		return
	}

	room := w.getRoom(roomUUID)
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
