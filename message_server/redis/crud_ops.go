package redis

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"

	"github.com/redis/go-redis/v9"
)

/*
Deletes a value by its key in the Redis database. Returns the number of
objects that were deleted. This function works for both single value and
multi-value keypairs, hence why `DelA` is not a valid function. If the key
doesn't exist, then this value will be 0. Applicable to D in CRUD.
*/
func Del(c *redis.Client, ctx context.Context, keys ...string) (int64, error) {
	return c.Del(ctx, keys...).Result()
}

/*
Gets the value for a key in the Redis database. If the key doesn't exist, then
`nil` will be emitted. Applicable to R in CRUD. See: https://stackoverflow.com/a/53697645
*/
func Get[T any](c *redis.Client, ctx context.Context, key string, dest *T) error {
	//Get from Redis
	p, err := c.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	//Unmarshal from GOB
	b := bytes.NewBuffer([]byte(p))
	dec := gob.NewDecoder(b)
	return dec.Decode(dest)
}

/*
Gets the array of values for a key in the Redis database. If the key
doesn't exist, then an empty array will be emitted. Applicable to R in
CRUD. See: https://stackoverflow.com/a/53697645
*/
func GetA[T any](c *redis.Client, ctx context.Context, key string) ([]T, error) {
	//Get from Redis
	ps, err := c.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	//Allocate space for the incoming elements
	dest := make([]T, len(ps))

	//Loop over each incoming object
	for i, p := range ps {
		//Unmarshal from GOB
		b := bytes.NewBuffer([]byte(p))
		dec := gob.NewDecoder(b)
		//if err := dec.Decode(&(dest)[i]); err != nil {
		if err := dec.Decode(&(dest)[i]); err != nil {
			return nil, err
		}
	}
	return dest, nil
}

/*
Sets a key and value in the Redis database. If the key already exists, its
value is updated. If not, then a new key is created. Applicable to C, U in
CRUD. See: https://stackoverflow.com/a/53697645
*/
func Set[T any](c *redis.Client, ctx context.Context, key string, value T) error {
	//Marshal to GOB
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(value); err != nil {
		return err
	}

	//Add to Redis
	return c.Set(ctx, key, b.Bytes(), time.Duration(0)).Err()
}
