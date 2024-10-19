package chat

import (
	"time"

	"wraith.me/message_server/util"
)

// Represents the last message in a chat room for quick preview.
type LastMessage struct {
	MessageID util.UUID `json:"message_id" bson:"message_id"`
	Content   string    `json:"content" bson:"content"`
}
