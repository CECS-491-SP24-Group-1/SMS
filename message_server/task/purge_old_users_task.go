package task

import (
	"context"
	"fmt"
	"time"

	
	"go.mongodb.org/mongo-driver/bson"
	"wraith.me/message_server/globals"
	

	"github.com/redis/go-redis/v9"
	"wraith.me/message_server/db" // Accessing the db from message_server
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
	//Query for inactive users in the MongoDB collection
	filter := bson.M{"should_purge": true}

	// Deleting the users flagged as inactive
	deleteResult, err := globals.UC.Collection.RemoveAll(ftt.CTX, filter)
	if err != nil {
		fmt.Printf("Error deleting inactive users: %v\n", err)
		return
	}

	// Log how many users were deleted
	fmt.Printf("Deleted %d inactive users.\n", deleteResult.DeletedCount)
	
}
