package mongoutil

import (
	"context"
	"fmt"

	"github.com/qiniu/qmgo"
	"wraith.me/message_server/pkg/db"
	chatroom "wraith.me/message_server/pkg/schema/chat_room"
	"wraith.me/message_server/pkg/util"
)

// Converts an aggregation result of `_id: UUID` to an array.
func Aggregate2IDArr(agg qmgo.AggregateI) ([]util.UUID, error) {
	//Perform the aggregation
	mp := make([]map[string]util.UUID, 0) //Array of single valued maps, keyed by `_id`
	err := agg.All(&mp)
	if err != nil {
		return nil, err
	}

	//Loop over the aggregation array and collect all UUIDs to an array of just UUIDs
	out := make([]util.UUID, len(mp))
	for i, uuid := range mp {
		out[i] = uuid["_id"]
	}
	return out, nil
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
