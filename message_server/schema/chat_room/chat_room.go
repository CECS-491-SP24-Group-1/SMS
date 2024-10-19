package chatroom

import (
	"time"

	"wraith.me/message_server/http_types/ws/chat"
	"wraith.me/message_server/util"
)

// Represents a chat room containing multiple participants and messages.
type ChatRoom struct {
	// Unique identifier for the chat room.
	ID util.UUID `json:"id" bson:"_id"`

	// The list of participants in the chat room, represented by their UUIDs.
	Participants []util.UUID `json:"participants" bson:"participants"`

	// The list of messages in the chat room.
	Messages []chat.ChatMessage `json:"messages" bson:"-"`

	// Summary of the last message in the chat room for quick access.
	LastMessage chat.LastMessage `json:"last_message" bson:"-"`

	// The timestamp of when the chat room was created.
	CreatedAt time.Time `json:"created_at" bson:"created_at"` //TODO: added by QMgoBase

	// The timestamp of when the chat room was last updated (when the last message was sent).
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"` //TODO: added by QMgoBase
}

// Adds a new message to the chat room and updates the last message.
func (c *ChatRoom) AddMessage(message chat.ChatMessage) {
	c.Messages = append(c.Messages, message)
	c.LastMessage = chat.LastMessage{
		ID:      message.ID,
		Content: message.Content,
	}
	c.UpdatedAt = message.ID.Time()
}
