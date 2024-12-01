package response

import (
	chatroom "wraith.me/message_server/pkg/schema/chat_room"
	"wraith.me/message_server/pkg/util"
)

// Represents a member of a chat room.
type RoomMember struct {
	//The ID of the member.
	ID util.UUID `json:"id"`

	//The username of the user.
	Username string `json:"username"`

	//The display name of the user.
	DisplayName string `json:"display_name"`

	//Whether this user is the currently logged in user.
	IsMe bool `json:"is_me"`

	//The role of the user.
	Role chatroom.Role `json:"role"`

	//Whether the user is currently online.
	IsOnline bool `json:"is_online"`
}
