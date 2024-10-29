package wschat

import (
	"sync"

	"github.com/olahol/melody"
	chatroom "wraith.me/message_server/pkg/schema/chat_room"
	"wraith.me/message_server/pkg/util"
)

// Represents an active chatroom instance.
type WSRoom struct {
	ID           util.UUID
	sessions     map[*melody.Session]*UserData
	userIDs      map[util.UUID]*melody.Session
	participants chatroom.MembershipList
	mu           sync.RWMutex
}

// Creates a new room.
func NewRoom(id util.UUID) *WSRoom {
	return &WSRoom{
		ID:           id,
		sessions:     make(map[*melody.Session]*UserData),
		userIDs:      make(map[util.UUID]*melody.Session),
		participants: make(chatroom.MembershipList),
	}
}

/*
Checks if a user ID already has a session and adds it if not.
Returns true if the session was added, false if the user already had a session.
*/
func (r *WSRoom) AddSession(s *melody.Session, userData *UserData) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if existingSession, exists := r.userIDs[userData.ID]; exists {
		// User already has a session
		if existingSession != s {
			// Remove the old session if it's different
			delete(r.sessions, existingSession)
			r.sessions[s] = userData
			r.userIDs[userData.ID] = s
		}
		return false
	}

	// Add new session
	r.sessions[s] = userData
	r.userIDs[userData.ID] = s
	return true
}

// Broadcasts a message to the room.
func (r *WSRoom) Broadcast(msg []byte, excludes ...*melody.Session) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	//Create a map for quick lookup of excluded sessions
	excludeMap := make(map[*melody.Session]bool)
	for _, exclude := range excludes {
		excludeMap[exclude] = true
	}

	//Iterate through the sessions map and send the message
	for session := range r.sessions {
		if !excludeMap[session] {
			session.Write(msg)
		}
	}
}

// Checks if a room has a session.
func (r *WSRoom) HasSession(s *melody.Session) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.sessions[s]
	return exists
}

// Checks if a user is in the room.
func (r *WSRoom) HasUser(uid util.UUID) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.userIDs[uid]
	return exists
}

// Checks if the room is empty.
func (r *WSRoom) IsEmpty() bool {
	return r.Size() == 0
}

// Removes a user from the room.
func (r *WSRoom) RemoveSession(s *melody.Session) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if userData, exists := r.sessions[s]; exists {
		delete(r.userIDs, userData.ID)
	}
	delete(r.sessions, s)
}

// Gets the number of active sessions in the room.
func (r *WSRoom) Size() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.sessions)
}

// Gets the user data associated with a session.
func (r *WSRoom) GetUserData(s *melody.Session) (*UserData, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	userData, exists := r.sessions[s]
	return userData, exists
}

// Gets all user data in the room.
func (r *WSRoom) GetAllUserData() []*UserData {
	r.mu.RLock()
	defer r.mu.RUnlock()
	userData := make([]*UserData, 0, len(r.sessions))
	for _, data := range r.sessions {
		userData = append(userData, data)
	}
	return userData
}

// Updates the user data for a given session.
func (r *WSRoom) UpdateUserData(s *melody.Session, newUserData *UserData) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if oldUserData, exists := r.sessions[s]; exists {
		delete(r.userIDs, oldUserData.ID)
		r.sessions[s] = newUserData
		r.userIDs[newUserData.ID] = s
	}
}
