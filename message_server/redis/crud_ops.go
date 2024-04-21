package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"wraith.me/message_server/db/mongoutil"
	"wraith.me/message_server/util"
)

const (
	//The default expiry of Redis tokens: no expiration.
	DEFAULT_EXPIRY = time.Duration(0)
)

//TODO: Add expiry methods too

/*
Creates a key and value in the Redis database. This function is an alias of
`Set()`. Applicable to C in CRUD. See: https://stackoverflow.com/a/53697645
*/
func Create[K uuid.UUID | mongoutil.UUID, V any](c *redis.Client, ctx context.Context, key K, value V) error {
	return Set(c, ctx, key, value)
}

/*
Creates a key and array of values in the Redis database. This function is an
alias of `SetA()`. Applicable to C in CRUD. See:
https://stackoverflow.com/a/53697645
*/
func CreateA[K uuid.UUID | mongoutil.UUID, V any](c *redis.Client, ctx context.Context, key K, values ...V) error {
	return SetA(c, ctx, key, values...)
}

/*
Creates a series of keys and object values in the Redis database from a map.
This function is an alias of `SetMany()`. Applicable to C in CRUD. See:
https://stackoverflow.com/a/53697645
*/
func CreateMany[K uuid.UUID | mongoutil.UUID, V any](c *redis.Client, ctx context.Context, kp map[K]V) error {
	return SetMany(c, ctx, kp)
}

/*
Creates a series of keys and object values in the Redis database from a map.
This function is an alias of `SetMany()`. Applicable to C in CRUD. See:
https://stackoverflow.com/a/53697645
*/
func CreateManyS[K uuid.UUID | mongoutil.UUID](c *redis.Client, ctx context.Context, kp map[K]string) error {
	return SetManyS(c, ctx, kp)
}

/*
Creates a key and array of strings in the Redis database. This function is an
alias of `SetSA()`. Applicable to C in CRUD. See:
https://stackoverflow.com/a/53697645
*/
func CreateSA[K uuid.UUID | mongoutil.UUID](c *redis.Client, ctx context.Context, key K, values ...string) error {
	return SetSA(c, ctx, key, values...)
}

/*
Deletes a value or array of values by its key in the Redis database.
Returns the number of objects that were deleted. This function works for
both single value and multi-value keypairs, hence why `DelA` is  not a
valid function. If the key doesn't exist, then this value will be 0.
Applicable to D in CRUD.
*/
func Del[K uuid.UUID | mongoutil.UUID](c *redis.Client, ctx context.Context, keys ...K) (int64, error) {
	skeys := make([]string, len(keys))
	for i, v := range keys {
		skeys[i] = u2s(v)
	}
	return c.Del(ctx, skeys...).Result()
}

/*
Gets the value for a key in the Redis database. If the key doesn't exist,
then `nil` will be emitted. Applicable to R in CRUD. See:
https://stackoverflow.com/a/53697645
*/
func Get[K uuid.UUID | mongoutil.UUID, V any](c *redis.Client, ctx context.Context, key K, dest *V) (err error) {
	//Get from Redis
	p, err := c.Get(ctx, u2s(key)).Result()
	if err != nil {
		return err
	}

	//Unmarshal from bytes
	*dest, err = util.FromGOBBytes[V]([]byte(p))
	return
}

/*
Gets the array of values for a key in the Redis database. If the key
doesn't exist, then an empty array will be emitted. Applicable to R in
CRUD. See: https://stackoverflow.com/a/53697645
*/
func GetA[K uuid.UUID | mongoutil.UUID, V any](c *redis.Client, ctx context.Context, key K) ([]V, error) {
	//Get from Redis
	ps, err := c.LRange(ctx, u2s(key), 0, -1).Result()
	if err != nil {
		return nil, err
	}

	//Allocate space for the incoming elements
	dest := make([]V, len(ps))

	//Loop over each incoming object
	for i, p := range ps {
		//Unmarshal from bytes
		obj, err := util.FromGOBBytes[V]([]byte(p))
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
func GetAt[K uuid.UUID | mongoutil.UUID, V any](c *redis.Client, ctx context.Context, key K, idx int64, dest *V) (err error) {
	//Get from Redis
	p, err := c.LIndex(ctx, u2s(key), idx).Result()
	if err != nil {
		return err
	}

	//Unmarshal from bytes
	*dest, err = util.FromGOBBytes[V]([]byte(p))
	return
}

/*
Gets an array of objects for an array of keys in the Redis database.
The values must be present in the database for this function to succeed.
Applicable to R in CRUD. See: https://stackoverflow.com/a/53697645
*/
func GetMany[K uuid.UUID | mongoutil.UUID, V any](c *redis.Client, ctx context.Context, keys ...K) ([]V, MultiRedisErr) {
	//Create the output array, matching the size of the input key array
	dest := make([]V, len(keys))

	//Create a Redis pipeline
	pl := c.TxPipeline()

	//Loop over the input key array and queue each value to be fetched from Redis
	for _, key := range keys {
		//Query Redis for the item via the pipeline
		if err := pl.Get(ctx, u2s(key)).Err(); err != nil {
			return nil, MultiRedisErr{fmt.Errorf("pipeline queue err; %v", err), []int{}}
		}
	}

	//Execute the commands in the pipeline and get the results array
	//Only bail if a non-nil error is not a redis nil error
	resl, err := pl.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, MultiRedisErr{fmt.Errorf("pipeline exec err; %v", err), []int{}}
	}

	//Loop over the fetched strings
	problematicIndices := []int{}
	for i, res := range resl {
		//Check if the current result is a cache miss (`redis.Nil`)
		if res.Err() != nil && res.Err() == redis.Nil {
			//fmt.Printf("[RCRUD_GetMany] Cache miss @ idx %d\n", i)
			//Simply add the index to the list of those that are problematic and skip the iteration
			problematicIndices = append(problematicIndices, i)
			continue
		}

		//Check if the current result is an error
		if res.Err() != nil {
			return nil, MultiRedisErr{fmt.Errorf("error for res #%d: %v", i+1, res.Err()), []int{i}}
		}

		//Type assert the result to a `StringCmd`
		sr, ok := res.(*redis.StringCmd)
		if !ok {
			return nil, MultiRedisErr{fmt.Errorf("string assert err for res #%d", i+1), []int{i}}
		}

		//Unmarshal the value string to the target type add it to the output array
		obj, err := util.FromGOBBytes[V]([]byte(sr.Val()))
		if err != nil {
			return nil, MultiRedisErr{fmt.Errorf("unmarshal err for res #%d: %v", i+1, err), []int{i}}
		}
		dest[i] = obj
	}

	//Return the list of items and a `Redis.Nil` error if there was at least one problematic index
	var oerr error = nil
	if len(problematicIndices) > 0 {
		oerr = redis.Nil
	}
	return dest, MultiRedisErr{oerr, problematicIndices}
}

/*
Gets an array of strings for an array of keys in the Redis database.
The values must be present in the database for this function to succeed.
Applicable to R in CRUD. See: https://stackoverflow.com/a/53697645
*/
func GetManyS[K uuid.UUID | mongoutil.UUID](c *redis.Client, ctx context.Context, keys ...K) ([]string, MultiRedisErr) {
	//Create the output array, matching the size of the input key array
	dest := make([]string, len(keys))

	//Create a Redis pipeline
	pl := c.TxPipeline()

	//Loop over the input key array and queue each value to be fetched from Redis
	for _, key := range keys {
		//Query Redis for the item via the pipeline
		if err := pl.Get(ctx, u2s(key)).Err(); err != nil {
			return nil, MultiRedisErr{fmt.Errorf("pipeline queue err; %v", err), []int{}}
		}
	}

	//Execute the commands in the pipeline and get the results array
	//Only bail if a non-nil error is not a redis nil error
	resl, err := pl.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, MultiRedisErr{fmt.Errorf("pipeline exec err; %v", err), []int{}}
	}

	//Loop over the fetched strings
	problematicIndices := []int{}
	for i, res := range resl {
		//Check if the current result is a cache miss (`redis.Nil`)
		if res.Err() != nil && res.Err() == redis.Nil {
			fmt.Printf("[RCRUD_GetMany] Cache miss @ idx %d\n", i)
			//Simply add the index to the list of those that are problematic and skip the iteration
			problematicIndices = append(problematicIndices, i)
			continue
		}

		//Check if the current result is an error
		if res.Err() != nil {
			return nil, MultiRedisErr{fmt.Errorf("error for res #%d: %v", i+1, res.Err()), []int{i}}
		}

		//Type assert the result to a `StringCmd`
		sr, ok := res.(*redis.StringCmd)
		if !ok {
			return nil, MultiRedisErr{fmt.Errorf("string assert err for res #%d", i+1), []int{i}}
		}

		//Add the string result to the output array
		dest[i] = sr.Val()
	}

	//Return the list of items and a `Redis.Nil` error if there was at least one problematic index
	var oerr error = nil
	if len(problematicIndices) > 0 {
		oerr = redis.Nil
	}
	return dest, MultiRedisErr{oerr, problematicIndices}
}

/*
Gets the array of strings for a key in the Redis database. If the key
doesn't exist, then an empty array will be emitted. Applicable to R in
CRUD. See: https://stackoverflow.com/a/53697645
*/
func GetSA[K uuid.UUID | mongoutil.UUID](c *redis.Client, ctx context.Context, key K) ([]string, error) {
	//Get from Redis
	ps, err := c.LRange(ctx, u2s(key), 0, -1).Result()
	if err != nil {
		return nil, err
	}

	//Return the string array
	return ps, nil
}

/*
Sets a key and value in the Redis database. If the key already exists, its
value is updated. If not, then a new key is created. Applicable to C, U in
CRUD. See: https://stackoverflow.com/a/53697645
*/
func Set[K uuid.UUID | mongoutil.UUID, V any](c *redis.Client, ctx context.Context, key K, value V) error {
	//Marshal to bytes
	bytes, err := util.ToGOBBytes(value)
	if err != nil {
		return err
	}

	//Add to Redis
	return c.Set(ctx, u2s(key), bytes, time.Duration(0)).Err()
}

/*
Sets a key and array of values in the Redis database. If the key already exists,
its old contents are discarded and its value array is replaced with this one.
Applicable to U in CRUD. See: https://stackoverflow.com/a/53697645
*/
func SetA[K uuid.UUID | mongoutil.UUID, V any](c *redis.Client, ctx context.Context, key K, values ...V) error {
	//Create a Redis pipeline
	pl := c.TxPipeline()

	//Check if the key exists and delete it if it does
	exists, err := c.Exists(ctx, u2s(key)).Result()
	if err == nil && exists > 0 {
		pl.Del(ctx, u2s(key))
	}

	//Loop over each incoming object
	for _, value := range values {
		//Marshal the current value to bytes
		bytes, err := util.ToGOBBytes(value)
		if err != nil {
			return err
		}

		//Add the item to Redis
		if err := pl.RPush(ctx, u2s(key), bytes).Err(); err != nil {
			return err
		}
	}

	//Execute the pipeline
	_, err = pl.Exec(ctx)
	return err
}

/*
Sets a new value of a specific array item for a key in the Redis database.
The array item must be present for this function to succeed, which is zero-
indexed. Applicable to U in CRUD. See: https://stackoverflow.com/a/53697645
*/
func SetAt[K uuid.UUID | mongoutil.UUID, V any](c *redis.Client, ctx context.Context, key K, idx int64, value V) error {
	//Marshal to bytes
	bytes, err := util.ToGOBBytes(value)
	if err != nil {
		return err
	}

	//Add to Redis
	return c.LSet(ctx, u2s(key), idx, bytes).Err()
}

/*
Sets a series of keys and object values in the Redis database from a map.
If the key already exists, its value is updated. If not, then a new key is
created. Applicable to C, U in CRUD. See: https://stackoverflow.com/a/53697645
*/
func SetMany[K uuid.UUID | mongoutil.UUID, V any](c *redis.Client, ctx context.Context, kp map[K]V) error {
	//Create a Redis pipeline
	pl := c.TxPipeline()

	//Loop over the input map and add the pairing to the pipeline
	for key, val := range kp {
		//Marshal the current value to a byte string
		bytes, err := util.ToGOBBytes(val)
		if err != nil {
			return err
		}

		//Add the current keypair to Redis
		if err := pl.Set(ctx, u2s(key), bytes, time.Duration(0)).Err(); err != nil {
			return err
		}
	}

	//Add the items to Redis by executing the pipeline
	_, err := pl.Exec(ctx)
	return err
}

/*
Sets a series of keys and string values in the Redis database from a map.
If the key already exists, its value is updated. If not, then a new key is
created. Applicable to C, U in CRUD. See: https://stackoverflow.com/a/53697645
*/
func SetManyS[K uuid.UUID | mongoutil.UUID](c *redis.Client, ctx context.Context, kp map[K]string) error {
	//Create a Redis pipeline
	pl := c.TxPipeline()

	//Loop over the input map and add the pairing to the pipeline
	for key, val := range kp {
		if err := pl.Set(ctx, u2s(key), val, time.Duration(0)).Err(); err != nil {
			return err
		}
	}

	//Add the items to Redis by executing the pipeline
	_, err := pl.Exec(ctx)
	return err
}

/*
Sets a key and array of strings in the Redis database. If the key already exists,
its old contents are discarded and its value array is replaced with this one.
Applicable to U in CRUD. See: https://stackoverflow.com/a/53697645
*/
func SetSA[K uuid.UUID | mongoutil.UUID](c *redis.Client, ctx context.Context, key K, values ...string) error {
	//Create a Redis pipeline
	pl := c.TxPipeline()

	//Check if the key exists and delete it if it does
	exists, err := c.Exists(ctx, u2s(key)).Result()
	if err == nil && exists > 0 {
		pl.Del(ctx, u2s(key))
	}

	//Add the items to Redis
	if err := pl.RPush(ctx, u2s(key), values).Err(); err != nil {
		return err
	}

	//Execute the pipeline
	_, err = pl.Exec(ctx)
	return err
}

// Converts a UUID-like object to a string
func u2s[T uuid.UUID | mongoutil.UUID](uuid T) string {
	return fmt.Sprintf("%s", uuid)
}
