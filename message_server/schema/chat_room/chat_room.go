package chat_room

import (
	"time"

	"wraith.me/message_server/obj"
)

// Represents a chat room containing multiple participants and messages.
type ChatRoom struct {
	obj.Identifiable `json:",inline" bson:",inline,squash"`

	// The list of participants in the chat room.
	Participants []obj.Participant `json:"participants" bson:"participants"`

	// The list of messages in the chat room.
	Messages []ChatMessage `json:"messages" bson:"messages"`

	// Summary of the last message in the chat room for quick access.
	LastMessage LastMessage `json:"last_message" bson:"last_message"`

	// The timestamp of when the chat room was created.
	CreatedAt time.Time `json:"created_at" bson:"created_at"`

	// The timestamp of when the chat room was last updated (when the last message was sent).
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// Adds a new message to the chat room and updates the last message.
func (c *ChatRoom) AddMessage(message ChatMessage) {
	c.Messages = append(c.Messages, message)
	c.LastMessage = LastMessage{
		MessageID: message.ID,
		Timestamp: message.Timestamp,
		Content:   message.Content,
	}
	c.UpdatedAt = message.Timestamp
}
