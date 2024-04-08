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
func rstructSet[T any](c *redis.Client, ctx context.Context, key string, value *T) error {
	//Marshal to GOB
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(value); err != nil {
		return err
	}

	//Add to Redis
	return c.Set(ctx, key, b.Bytes(), time.Duration(0)).Err()
}

// See: https://stackoverflow.com/a/53697645
func rstructGet[T any](c *redis.Client, ctx context.Context, key string, dest *T) error {
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

func TestRedisObjPush(t *testing.T) {
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

	//Connect to Redis
	red := redisInit()

	//Push the objects into the Redis database
	if err := rstructSet(red, context.Background(), foo1.ID.String(), &foo1); err != nil {
		t.Error(err)
	}
	if err := rstructSet(red, context.Background(), foo2.ID.String(), &foo2); err != nil {
		t.Error(err)
	}
	if err := rstructSet(red, context.Background(), foo3.ID.String(), &foo3); err != nil {
		t.Error(err)
	}

	//Get the objects from Redis
	var rfoo1 Foo
	if err := rstructGet(red, context.Background(), foo1.ID.String(), &rfoo1); err != nil {
		t.Error(err)
	}
	if !eq(&foo1, &rfoo1) {
		fmt.Printf("foo1 : %v\n", foo1)
		fmt.Printf("rfoo1: %v\n", rfoo1)
		t.Errorf("foo1 != rfoo1")
	}
	var rfoo2 Foo
	if err := rstructGet(red, context.Background(), foo2.ID.String(), &rfoo2); err != nil {
		t.Error(err)
	}
	if !eq(&foo2, &rfoo2) {
		fmt.Printf("foo2 : %v\n", foo2)
		fmt.Printf("rfoo2: %v\n", rfoo2)
		t.Errorf("foo2 != rfoo2")
	}
	var rfoo3 Foo
	if err := rstructGet(red, context.Background(), foo3.ID.String(), &rfoo3); err != nil {
		t.Error(err)
	}
	if !eq(&foo3, &rfoo3) {
		fmt.Printf("foo3 : %v\n", foo3)
		fmt.Printf("rfoo3: %v\n", rfoo3)
		t.Errorf("foo3 != rfoo3")
	}
}
