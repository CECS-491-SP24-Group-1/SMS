package task

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Purges old users from the database; implements `Task`.
type PurgeOldUsersTask struct {
	//Defines the duration between runs.
	TQ time.Duration

	//The Redis client to use in transactions.
	RC *redis.Client

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
}
