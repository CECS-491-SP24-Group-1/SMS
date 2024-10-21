package wschat

import (
	"sync"

	"github.com/olahol/melody"
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
	rooms  map[util.UUID]*Room
	roomMu sync.RWMutex
}

// Gets the currently active chat server instance.
func GetInstance() *Server {
	once.Do(func() {
		instance = &Server{
			melody: melody.New(),
			mutex:  &sync.Mutex{},
			rooms:  make(map[util.UUID]*Room),
		}
		instance.setupHandlers()
	})
	return instance
}

// Gets the backend Melody handler for the server.
func (w *Server) GetMelody() *melody.Melody {
	return w.melody
}

func (w *Server) getOrCreateRoom(id util.UUID) *Room {
	w.roomMu.Lock()
	defer w.roomMu.Unlock()

	if room, exists := w.rooms[id]; exists {
		return room
	}

	room := NewRoom(id)
	w.rooms[id] = room
	return room
}

func (w *Server) getRoom(id util.UUID) *Room {
	w.roomMu.RLock()
	defer w.roomMu.RUnlock()

	return w.rooms[id]
}

func (w *Server) removeRoom(id util.UUID) {
	w.roomMu.Lock()
	defer w.roomMu.Unlock()

	delete(w.rooms, id)
}

func isValidMessage(msg []byte) bool {
	return len(msg) > 0
}
