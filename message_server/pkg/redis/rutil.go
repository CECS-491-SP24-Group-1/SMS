package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

/*
Checks if a Redis database has a given key. Ensure Redis is accessible and
no errors occur. This function assumes no error will occur. If one does occur,
then it'll silently fail and return false.
*/
func RHas(red *redis.Client, key string, ctx context.Context) bool {
	_, err := red.Get(ctx, key).Result()
	return err != nil //`err` can be anything; this is a silent error
}
