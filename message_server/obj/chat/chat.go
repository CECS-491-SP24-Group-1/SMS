package chat

import (
	"time"

	"wraith.me/message_server/util"
)

// Represents an individual message within a chat room.
type ChatMessage struct {

	// The ID of the user who sent the message.
	SenderID util.UUID `json:"sender_id" bson:"sender_id"`

	// The content of the message.
	Content string `json:"content" bson:"content"`

	// The timestamp when the message was sent.
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`

	// Status of the message (e.g., sent, delivered, read).
	Status string `json:"status" bson:"status"`

	// Attachments within the message (if any).
	//TODO
}
