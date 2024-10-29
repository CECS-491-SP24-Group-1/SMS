package wschat

import (
	"sync"

	"github.com/olahol/melody"
	chatroom "wraith.me/message_server/pkg/schema/chat_room"
	"wraith.me/message_server/pkg/util"
)

var (
	// Holds the instance object for the global WebSocket chats server.
	instance *Server

	// Guard mutex to ensure that only one singleton object is created.
	once sync.Once
)

//-- SINGLETON: Server

/*
Represents a Melody WebSocket chat server. This struct acts as a singleton
wrapper on a Melody server that's responsible for chats.
*/
type Server struct {
	melody *melody.Melody
	mutex  *sync.Mutex
	rooms  map[util.UUID]*WSRoom
	roomMu sync.RWMutex
}

// Gets the currently active chat server instance.
func GetInstance() *Server {
	once.Do(func() {
		instance = &Server{
			melody: melody.New(),
			mutex:  &sync.Mutex{},
			rooms:  make(map[util.UUID]*WSRoom),
		}
		instance.setupHandlers()
	})
	return instance
}

// Gets the backend Melody handler for the server.
func (w *Server) GetMelody() *melody.Melody {
	return w.melody
}

// Attempts to get a room by ID. If it doesn't exist, a new one is created.
func (w *Server) getOrCreateRoom(id util.UUID, participants chatroom.MembershipList) *WSRoom {
	w.roomMu.Lock()
	defer w.roomMu.Unlock()

	if room, exists := w.rooms[id]; exists {
		return room
	}

	room := NewRoom(id)
	room.participants = participants
	w.rooms[id] = room
	return room
}

// Gets a room by ID.
func (w *Server) getRoom(id util.UUID) *WSRoom {
	w.roomMu.RLock()
	defer w.roomMu.RUnlock()

	return w.rooms[id]
}

// Removes a room from the handler.
func (w *Server) removeRoom(id util.UUID) {
	w.roomMu.Lock()
	defer w.roomMu.Unlock()

	delete(w.rooms, id)
}

// Sets up the handlers for the chat server.
func (w *Server) setupHandlers() {
	w.melody.HandleConnect(w.handleConnect)
	w.melody.HandleDisconnect(w.handleDisconnect)
	w.melody.HandleMessage(w.handleMessage)
}
