package globals

import (
	"github.com/redis/go-redis/v9"
)

// Rcl is the global Redis client that can be accessed from anywhere in the application
var Rcl *redis.Client
