package wschat

import (
	"sync"

	"github.com/olahol/melody"
	"wraith.me/message_server/pkg/util"
)

// Represents an active chatroom instance.
type Room struct {
	ID       util.UUID
	sessions map[*melody.Session]bool
	mu       sync.RWMutex
}

// Creates a new room.
func NewRoom(id util.UUID) *Room {
	return &Room{
		ID:       id,
		sessions: make(map[*melody.Session]bool),
	}
}

// Adds a user to the room.
func (r *Room) AddSession(s *melody.Session) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[s] = true
}

// Broadcasts a message to the room.
func (r *Room) Broadcast(msg []byte, exclude *melody.Session) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	//TODO: add ability to exclude multiple sessions
	for session := range r.sessions {
		if session != exclude {
			session.Write(msg)
		}
	}
}

// Checks if a user is in the room.
func (r *Room) HasSession(s *melody.Session) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.sessions[s]
	return exists
}

// Checks if the room is empty.
func (r *Room) IsEmpty() bool {
	return r.Size() == 0
}

// Removes a user from the room.
func (r *Room) RemoveSession(s *melody.Session) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.sessions, s)
}

// Gets the number of active sessions in the room.
func (r *Room) Size() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.sessions)
}
