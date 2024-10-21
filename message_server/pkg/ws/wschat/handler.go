package wschat

import (
	"github.com/olahol/melody"
	"wraith.me/message_server/pkg/util"
)

// Sets up the handlers for the chat server,
func (w *Server) setupHandlers() {
	w.melody.HandleConnect(w.handleConnect)
	w.melody.HandleDisconnect(w.handleDisconnect)
	w.melody.HandleMessage(w.handleMessage)
}

// Handles new connections to the chat server.
func (w *Server) handleConnect(s *melody.Session) {
	roomID, ok := s.Request.Context().Value(WSChatCtxRoomIDKey).(util.UUID)
	if !ok {
		s.Write([]byte("Invalid room ID"))
		s.Close()
		return
	}

	room := w.getOrCreateRoom(roomID)
	room.AddSession(s)

	s.Write([]byte("Welcome to the room!"))
	room.Broadcast([]byte("A user has connected to the room"), s)
}

// Handles disconnections from the chat server.
func (w *Server) handleDisconnect(s *melody.Session) {
	roomUUID, ok := s.Request.Context().Value(WSChatCtxRoomIDKey).(util.UUID)
	if !ok {
		return
	}

	room := w.getRoom(roomUUID)
	if room != nil {
		room.RemoveSession(s)
		room.Broadcast([]byte("A user has disconnected from the room"), s)
		if room.IsEmpty() {
			w.removeRoom(roomUUID)
		}
	}
}

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
