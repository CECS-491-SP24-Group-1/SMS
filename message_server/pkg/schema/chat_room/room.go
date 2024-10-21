package chatroom

import (
	"wraith.me/message_server/pkg/db"
	"wraith.me/message_server/pkg/util"
)

// Represents a chat room containing multiple participants and messages.
type Room struct {
	db.DBObj `bson:",inline"`

	// Unique identifier for the chat room.
	ID util.UUID `json:"id" bson:"_id"`

	// The list of participants in the chat room, represented by their UUIDs.
	Participants map[util.UUID]Role `json:"participants" bson:"participants"`

	// The list of messages in the chat room.
	//Messages []chat.ChatMessage `json:"messages" bson:"-"`

	// Summary of the last message in the chat room for quick access.
	//LastMessage chat.LastMessage `json:"last_message" bson:"-"`
}

// Creates a new chat room.
func NewRoom(owner util.UUID, participants ...util.UUID) Room {
	//Create the map of members and add to it
	members := make(map[util.UUID]Role)
	for _, participant := range participants {
		members[participant] = RoleMEMBER
	}

	//Create the room
	return Room{
		DBObj:        db.NewDBObj(),
		ID:           util.MustNewUUID7(),
		Participants: members,
	}
}

/*
// Adds a new message to the chat room and updates the last message.
func (c *Room) AddMessage(message chat.ChatMessage) {
	c.Messages = append(c.Messages, message)
	c.LastMessage = chat.LastMessage{
		ID:      message.ID,
		Content: message.Content,
	}
	//c.UpdatedAt = message.Timestamp
}

// CreateChatRoom inserts a new chat room into the MongoDB collection.
func CreateChatRoom(chatRoom *chatroom.ChatRoom) error {
	client := db.GetInstance().GetClient() // Call GetInstance() from singleton.go
	if client == nil {
		return fmt.Errorf("MongoDB client is not initialized")
	}

	collection := client.Database("your_database_name").Collection("chat_rooms")

	// Insert the chat room document into the collection
	_, err := collection.InsertOne(context.TODO(), chatRoom)
	return err
}
*/
