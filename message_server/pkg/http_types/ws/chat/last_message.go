package chat

import (
	"wraith.me/message_server/pkg/util"
)

//May be unnecessary

// Represents the last message in a chat room for quick preview.
type LastMessage struct {
	ID      util.UUID `json:"id" bson:"_id"`
	Content string    `json:"content" bson:"content"`
}
