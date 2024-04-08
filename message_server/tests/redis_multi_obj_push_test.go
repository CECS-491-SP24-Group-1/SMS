package tests

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// See: https://stackoverflow.com/a/53697645
func rstructSetA[T any](c *redis.Client, ctx context.Context, key string, values []*T) error {
	//Create the output byte array
	bytea := make([][]byte, len(values))

	//Loop over each incoming object
	for i, value := range values {
		//Marshal to GOB
		var b bytes.Buffer
		enc := gob.NewEncoder(&b)
		if err := enc.Encode(value); err != nil {
			return err
		}
		bytea[i] = b.Bytes()
	}

	//Add each item to Redis
	for _, val := range bytea {
		if err := c.RPush(ctx, key, val).Err(); err != nil {
			return err
		}
	}
	return nil
}

// See: https://stackoverflow.com/a/53697645
func rstructGetA[T any](c *redis.Client, ctx context.Context, key string) ([]*T, error) {
	//Get from Redis
	ps, err := c.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	//Allocate space for the incoming elements
	dest := make([]*T, len(ps))

	//Loop over each incoming object
	for i, p := range ps {
		//Unmarshal from GOB
		b := bytes.NewBuffer([]byte(p))
		dec := gob.NewDecoder(b)
		if err := dec.Decode(&(dest)[i]); err != nil {
			return nil, err
		}
	}
	return dest, nil
}

func TestRedisMultiObjPush(t *testing.T) {
	//Define the object
	type Foo struct {
		ID            uuid.UUID
		Name          string
		Birthday      time.Time
		FavoriteFoods []string
	}

	eq := func(a *Foo, b *Foo) bool {
		return a.ID == b.ID && a.Name == b.Name && a.Birthday == b.Birthday && slices.Equal(a.FavoriteFoods, b.FavoriteFoods)
	}
	eqa := func(a []*Foo, b []*Foo) bool {
		return slices.EqualFunc(a, b, eq)
	}

	//Create some instances
	foo1 := Foo{
		ID:            uuid.New(),
		Name:          "John Doe",
		Birthday:      time.Now().Round(0),
		FavoriteFoods: []string{"carrots", "apples", "pasta"},
	}
	foo2 := Foo{
		ID:            uuid.New(),
		Name:          "Jane Doe",
		Birthday:      time.Now().Round(0),
		FavoriteFoods: []string{"bananas", "melons", "ice-cream"},
	}
	foo3 := Foo{
		ID:            uuid.New(),
		Name:          "Jin Doe",
		Birthday:      time.Now().Round(0),
		FavoriteFoods: []string{"ramen", "rice", "sushi"},
	}

	//Construct the array and an ID for it
	aid := uuid.New()
	objs := []*Foo{&foo1, &foo2, &foo3}

	//Connect to Redis
	red := redisInit()

	//Push the object list into the Redis database
	if err := rstructSetA(red, context.Background(), aid.String(), objs); err != nil {
		t.Error(err)
	}
	//Get the objects from Redis
	robjs, err := rstructGetA[Foo](red, context.Background(), aid.String())
	if err != nil {
		t.Error(err)
	}
	if !eqa(objs, robjs) {
		fmt.Printf("objs : %v\n", objs)
		fmt.Printf("robjs: %v\n", robjs)
		t.Errorf("objs != robjs")
	}
}
