package redis

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

/*
Creates a key and value in the Redis database. This function is an alias of
`Set()`. Applicable to C in CRUD. See: https://stackoverflow.com/a/53697645
*/
func Create[T any](c *redis.Client, ctx context.Context, key uuid.UUID, value T) error {
	return Set(c, ctx, key, value)
}

/*
Creates a key and array of values in the Redis database. This function is an
alias of `SetA()`. Applicable to C in CRUD. See:
https://stackoverflow.com/a/53697645
*/
func CreateA[T any](c *redis.Client, ctx context.Context, key uuid.UUID, values []T) error {
	return SetA(c, ctx, key, values)
}

/*
Deletes a value or array of values by its key in the Redis database.
Returns the number of objects that were deleted. This function works for
both single value and multi-value keypairs, hence why `DelA` is  not a
valid function. If the key doesn't exist, then this value will be 0.
Applicable to D in CRUD.
*/
func Del(c *redis.Client, ctx context.Context, keys ...uuid.UUID) (int64, error) {
	skeys := make([]string, len(keys))
	for i, v := range keys {
		skeys[i] = v.String()
	}
	return c.Del(ctx, skeys...).Result()
}

/*
Gets the value for a key in the Redis database. If the key doesn't exist,
then `nil` will be emitted. Applicable to R in CRUD. See:
https://stackoverflow.com/a/53697645
*/
func Get[T any](c *redis.Client, ctx context.Context, key uuid.UUID, dest *T) (err error) {
	//Get from Redis
	p, err := c.Get(ctx, key.String()).Result()
	if err != nil {
		return err
	}

	//Unmarshal from bytes
	*dest, err = fromBytes[T]([]byte(p))
	return
}

/*
Gets the array of values for a key in the Redis database. If the key
doesn't exist, then an empty array will be emitted. Applicable to R in
CRUD. See: https://stackoverflow.com/a/53697645
*/
func GetA[T any](c *redis.Client, ctx context.Context, key uuid.UUID) ([]T, error) {
	//Get from Redis
	ps, err := c.LRange(ctx, key.String(), 0, -1).Result()
	if err != nil {
		return nil, err
	}

	//Allocate space for the incoming elements
	dest := make([]T, len(ps))

	//Loop over each incoming object
	for i, p := range ps {
		//Unmarshal from bytes
		obj, err := fromBytes[T]([]byte(p))
		if err != nil {
			return nil, err
		}
		dest[i] = obj
	}
	return dest, nil
}

/*
Gets the value of a specific array item for a key in the Redis database.
The array item must be present for this function to succeed, which is zero-
indexed. Applicable to R in CRUD. See: https://stackoverflow.com/a/53697645
*/
func GetAt[T any](c *redis.Client, ctx context.Context, key uuid.UUID, idx int64, dest *T) (err error) {
	//Get from Redis
	p, err := c.LIndex(ctx, key.String(), idx).Result()
	if err != nil {
		return err
	}

	//Unmarshal from bytes
	*dest, err = fromBytes[T]([]byte(p))
	return
}

/*
Sets a key and value in the Redis database. If the key already exists, its
value is updated. If not, then a new key is created. Applicable to C, U in
CRUD. See: https://stackoverflow.com/a/53697645
*/
func Set[T any](c *redis.Client, ctx context.Context, key uuid.UUID, value T) error {
	//Marshal to bytes
	bytes, err := toBytes(value)
	if err != nil {
		return err
	}

	//Add to Redis
	return c.Set(ctx, key.String(), bytes, time.Duration(0)).Err()
}

/*
Sets a key and array of values in the Redis database. If the key already exists,
its old contents are discarded and its value array is replaced with this one.
Applicable to U in CRUD. See: https://stackoverflow.com/a/53697645
*/
func SetA[T any](c *redis.Client, ctx context.Context, key uuid.UUID, values []T) error {
	//Create a Redis pipeline
	pl := c.Pipeline()

	//Check if the key exists and delete it if it does
	exists, err := c.Exists(ctx, key.String()).Result()
	if err == nil && exists > 0 {
		pl.Del(ctx, key.String())
	}

	//Loop over each incoming object
	for _, value := range values {
		//Marshal the value to bytes
		bytes, err := toBytes(value)
		if err != nil {
			return err
		}

		//Add the item to Redis
		if err := pl.RPush(ctx, key.String(), bytes).Err(); err != nil {
			return err
		}
	}

	//Execute the pipeline
	_, err = pl.Exec(context.Background())
	return err
}

/*
Sets a new value of a specific array item for a key in the Redis database.
The array item must be present for this function to succeed, which is zero-
indexed. Applicable to U in CRUD. See: https://stackoverflow.com/a/53697645
*/
func SetAt[T any](c *redis.Client, ctx context.Context, key uuid.UUID, idx int64, value T) error {
	//Marshal to bytes
	bytes, err := toBytes(value)
	if err != nil {
		return err
	}

	//Add to Redis
	return c.LSet(ctx, key.String(), idx, bytes).Err()
}

// Unmarshals an object from a GOB byte stream.
func fromBytes[T any](target []byte) (T, error) {
	//Unmarshal from GOB
	var dest T
	b := bytes.NewBuffer([]byte(target))
	dec := gob.NewDecoder(b)
	err := dec.Decode(&dest)
	return dest, err
}

// Marshals an object to a GOB byte stream.
func toBytes[T any](target T) ([]byte, error) {
	//Marshal to GOB
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(target); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
