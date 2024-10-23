package wschat

import "wraith.me/message_server/pkg/util"

// Represents an object passed as context to the ws server from the HTTP handler.
type Context struct {
	RoomID   util.UUID //The ID the room.
	MemberID util.UUID //The ID of the user trying to connect.
}
