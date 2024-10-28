package chat

import (
	"encoding/json"

	"wraith.me/message_server/pkg/util"
)

var (
	// The prefix of a server command.
	CommandPrefix = '/'
)

// Represents an individual message within a chat room.
type Message struct {
	// Unique identifier for the message.
	ID util.UUID `json:"id"`

	// The type of the chat message, e.g., MESSAGE, EVENT.
	Type Type `json:"type"`

	// The ID of the user who sent the message.
	Sender util.UUID `json:"sender_id"`

	// The ID of the user who received the message.
	Recipient util.UUID `json:"recipient_id"`

	// The content of the message.
	Content string `json:"content"`

	// Status of the message (e.g., sent, delivered, read).
	//Status string `json:"status" bson:"status"`

	// Attachments within the message (if any).
	//Attachments []obj.Attachment `json:"attachments" bson:"attachments"`
}

// Creates a new message, with type field.
func NewMessageTyp(content string, sid util.UUID, rid util.UUID, typ Type) Message {
	return Message{
		ID:        util.MustNewUUID7(),
		Type:      typ,
		Sender:    sid,
		Recipient: rid,
		Content:   content,
	}
}

// Creates a new message.
func NewMessage(content string, sid util.UUID, rid util.UUID) Message {
	return NewMessageTyp(content, sid, rid, TypeSMSG)
}

// Checks whether this message is a command.
func (m Message) IsCommand() bool {
	return len(m.Content) > 0 && m.Content[0] == byte(CommandPrefix)
}

// Gets the command contained in the message, if it is one.
func (m Message) GetCommand() string {
	if !m.IsCommand() || len(m.Content) <= 1 {
		return ""
	} else {
		return m.Content[2:]
	}
}

// Marshals the message to JSON.
func (m Message) JSON() []byte {
	jsons, err := json.Marshal(m)
	if err != nil {
		panic("Message::JSON: " + err.Error())
	}
	return jsons
}
