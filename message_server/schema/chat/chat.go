package chat

import (
	"time"
	"wraith.me/message_server/util"
)

//-- CLASS: Message
//

// Represents a single message in a chat.
type Message struct {
	// Unique identifier for the message
	MessageID util.UUID `json:"message_id" bson:"message_id"`

	// ID of the sender of the message
	SenderID util.UUID `json:"sender_id" bson:"sender_id"`

	// The timestamp when the message was sent
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`

	// The content of the message
	Content string `json:"content" bson:"content"`

	// Status of the message (e.g., sent, delivered, read)
	Status string `json:"status" bson:"status"`

	// Attachments for the message (if any)
	Attachments []Attachment `json:"attachments" bson:"attachments"`
}

//-- CLASS: Attachment
//

// Represents an attachment in a message. Can be expanded when we add attachmentlist
// First pass, still needs work
type Attachment struct {
	// The type of attachment (e.g., image, file, video)
	Type string `json:"type" bson:"type"`

	// The URL where the attachment is stored
	URL string `json:"url" bson:"url"`

	// The filename of the attachment
	Filename string `json:"filename" bson:"filename"`
}

//-- CLASS: Chat
//

// Represents a chat conversation between participants.
type Chat struct {
	// Unique identifier for the chat
	ChatID util.UUID `json:"chat_id" bson:"chat_id"`

	// List of participants in the chat
	Participants []Participant `json:"participants" bson:"participants"`

	// List of messages in the chat
	Messages []Message `json:"messages" bson:"messages"`

	// Last message in the chat for quick reference
	LastMessage LastMessage `json:"last_message" bson:"last_message"`
}

//-- CLASS: Participant
//

// Represents a participant in a chat.
type Participant struct {
	// Unique identifier for the participant
	UserID util.UUID `json:"user_id" bson:"user_id"`

	// The username of the participant
	Username string `json:"username" bson:"username"`

}

//-- CLASS: LastMessage
//

// Represents the last message in a chat for preview purposes.
type LastMessage struct {
	// The ID of the last message
	MessageID util.UUID `json:"message_id" bson:"message_id"`

	// The timestamp of the last message
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`

	// The content of the last message
	Content string `json:"content" bson:"content"`
}

//-- Constructors

// NewChat creates a new chat object.
func NewChat(chatID util.UUID, participants []Participant) *Chat {
	return &Chat{
		ChatID:      chatID,
		Participants: participants,
		Messages:    make([]Message, 0),
	}
}

// NewMessage creates a new message object.
func NewMessage(messageID, senderID util.UUID, content string, status string, timestamp time.Time) *Message {
	return &Message{
		MessageID:   messageID,
		SenderID:    senderID,
		Content:     content,
		Status:      status,
		Timestamp:   timestamp,
		Attachments: make([]Attachment, 0),
	}
}

// AddMessage adds a message to a chat.
func (c *Chat) AddMessage(message Message) {
	c.Messages = append(c.Messages, message)
	c.LastMessage = LastMessage{
		MessageID: message.MessageID,
		Timestamp: message.Timestamp,
		Content:   message.Content,
	}
}


// AddAttachment adds an attachment to a message.
func (m *Message) AddAttachment(attachment Attachment) {
	m.Attachments = append(m.Attachments, attachment)
}
