

import (
	"time"

	"wraith.me/message_server/util"
	
)

// Represents an individual message within a chat room.
type ChatMessage struct {
	// Unique identifier for the message.
	ID util.UUID `json:"id" bson:"_id"`

	// The type of the chat message, e.g., MESSAGE, EVENT.
	Type chat.ChatType `json:"type" bson:"type"`

	// The ID of the user who sent the message.
	Sender util.UUID `json:"sender_id" bson:"sender_id"`

	// The ID of the user who received the message.
	Recipient util.UUID `json:"recipient_id" bson:"recipient_id"`

	// The content of the message.
	Content string `json:"content" bson:"content"`

	
	// Status of the message (e.g., sent, delivered, read).
	//Status string `json:"status" bson:"status"`

	// Attachments within the message (if any).
	//TODO
}
