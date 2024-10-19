package ws

import (
	"fmt"
	"sync"

	"github.com/olahol/melody"
)

//-- SINGLETON: WServer

/*
Represents a Melody WebSocket server. This struct acts as a singleton wrapper
on a Melody server.
*/
type WServer struct {
	melody *melody.Melody
	// Guard mutex to ensure atomicity during connect/disconnect operations.
	mutex *sync.Mutex
}

// Holds the instance object for the global WebSocket server.
var instance *WServer

// Guard mutex to ensure that only one singleton object is created.
var once sync.Once

// Gets the currently active Melody server instance.
func GetInstance() *WServer {
	once.Do(func() {
		instance = &WServer{
			melody: melody.New(),
			mutex:  &sync.Mutex{},
		}
	})
	return instance
}

/*
Gets the underlying Melody instance used to handle WebSocket connections.
*/
func (w *WServer) GetMelody() *melody.Melody {
	return w.melody
}

/*
Sends a message to all connected clients. Returns an error if sending fails.
*/
func (w *WServer) Broadcast(message []byte) error {
	if w.melody == nil {
		return fmt.Errorf("websocket: cannot broadcast; melody instance is nil")
	}

	return w.melody.Broadcast(message)
}
