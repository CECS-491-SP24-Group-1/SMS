package chatroom

import (
	"fmt"
	"math/big"

	"crypto/rand"

	"golang.org/x/exp/maps"
	"wraith.me/message_server/pkg/db"
	"wraith.me/message_server/pkg/util"
)

// A map that pairs a user ID with a role.
type MembershipList map[util.UUID]Role

// Represents a chat room containing multiple participants and messages.
type Room struct {
	db.DBObj `bson:",inline"`

	// Unique identifier for the chat room.
	ID util.UUID `json:"id" bson:"_id"`

	// The list of participants in the chat room, represented by their UUIDs.
	Participants MembershipList `json:"participants" bson:"participants"`

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

	//Add the owner as a participant
	members[owner] = RoleOWNER

	//Create the room
	return Room{
		DBObj:        db.NewDBObj(),
		ID:           util.MustNewUUID7(),
		Participants: members,
	}
}

// Adds a member to the chat room.
func (r *Room) AddMember(participant util.UUID) {
	//Add the user to the room and give them the member role
	r.Participants[participant] = RoleMEMBER
}

// Checks if a user is in the chat room.
func (r *Room) HasMember(participant util.UUID) bool {
	_, ok := r.Participants[participant]
	return ok
}

// Returns whether the room is empty, and thus, safe to remove.
func (r Room) IsEmpty() bool {
	return r.Size() < 1
}

// Removes a member from the room.
func (r *Room) RemoveMember(participant util.UUID) {
	//Get the role of the user to remove and the number of people in the room
	targetRole := r.Participants[participant]
	countBefore := len(r.Participants)

	//Remove the user from the room
	delete(r.Participants, participant)

	//Break if the room is now empty
	countAfter := len(r.Participants)
	if countAfter == 0 {
		return
	}

	//If the count is now different, then continue
	if countAfter < countBefore {
		//If the user that left was the owner, then pick a random user to be the owner
		if targetRole == RoleOWNER {
			//Get the list of users still in the room
			users := r.Users()

			//Generate a secure random number between 0 and max-1
			n, err := rand.Int(rand.Reader, big.NewInt(int64(len(users))))
			if err != nil {
				panic(fmt.Sprintf("RemoveMember::pickAWinner: %s", err))
			}

			//Pick a winner
			newOwner := users[n.Int64()]

			//Reassign the role of the picked user
			r.Participants[newOwner] = RoleOWNER
		}
	}
}

// Gets the number of users in the room.
func (r Room) Size() int {
	return len(r.Users())
}

// Gets the list of users currently in the room.
func (r Room) Users() []util.UUID {
	return maps.Keys(r.Participants)
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
