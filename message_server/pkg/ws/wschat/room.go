package wschat

import (
	"sync"

	"github.com/olahol/melody"
	"wraith.me/message_server/pkg/util"
)

type Room struct {
	ID       util.UUID
	sessions map[*melody.Session]bool
	mu       sync.RWMutex
}

func NewRoom(id util.UUID) *Room {
	return &Room{
		ID:       id,
		sessions: make(map[*melody.Session]bool),
	}
}

func (r *Room) AddSession(s *melody.Session) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[s] = true
}

func (r *Room) RemoveSession(s *melody.Session) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.sessions, s)
}

func (r *Room) HasSession(s *melody.Session) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.sessions[s]
	return exists
}

func (r *Room) IsEmpty() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.sessions) == 0
}

func (r *Room) Broadcast(msg []byte, exclude *melody.Session) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for session := range r.sessions {
		if session != exclude {
			session.Write(msg)
		}
	}
}
