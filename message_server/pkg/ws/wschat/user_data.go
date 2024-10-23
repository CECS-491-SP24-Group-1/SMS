package wschat

import (
	"wraith.me/message_server/pkg/util"
)

// UserData represents the information associated with a user in a room.
type UserData struct {
	ID util.UUID
	// Add more fields as needed, for example:
	// Name	string
	// Role	string
}
