package task

import (
	"context"
	"fmt"
	"time"
	"log"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"wraith.me/message_server/db"	// Accessing the db from message_server
)

// Purges old users from the database; implements `Task`.
type PurgeOldUsersTask struct {
	//Defines the duration between runs.
	TQ time.Duration

	//The Redis client to use in transactions.
	RC *redis.Client
	// The MongoDB client to use in transactions.
	MongoClient *db.MClient


	//The context to run the operations in.
	CTX context.Context
}

var _ Task = (*PurgeOldUsersTask)(nil) // Type assertion check to ensure compliance with `Task` interface.

func (ftt PurgeOldUsersTask) periodicDuration() time.Duration {
	return ftt.TQ
}

func (ftt PurgeOldUsersTask) runOnStart() {
	ftt.runPeriodically()
}

func (ftt PurgeOldUsersTask) runOnStop() {
	ftt.runPeriodically()
}

func (ftt PurgeOldUsersTask) runPeriodically() {
	fmt.Printf("fit tea task; run periodically; time: %s\n", time.Now().Format("15:04:05.000"))

	// Check the database connection
	if !ftt.MongoClient.IsConnected(){
		log.Println("MongoDB client is not connected, skipping purge task")
		return
	}

	// Accessing user collection
	user_collection := ftt.MongoClient.GetClient()

	// User flag filter
	filter := bson.M{
		"flags.should_purge": true,
		"flags.purge_by": bson.M{"$lt": time.Now()},
	}
	// Perform the deletion
	result, err := user_collection.RemoveAll(ftt.CTX, filter)
	if err != nil {
		log.Printf("Error purging users: %v", err)
		return
	}

	log.Printf("Purged %d users from the database", result.DeletedCount)
}
