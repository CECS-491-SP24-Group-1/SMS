package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"wraith.me/message_server/db/mongoutil"
	cr "wraith.me/message_server/redis"
)

func TestRedisKVSetS(t *testing.T) {
	//Set the number of items to insert
	n := 5

	//Get a Redis client
	r := redisInit()

	//Create a map of test strings and IDs to go with them
	kvs := make(map[uuid.UUID]string)
	keys := make([]uuid.UUID, n)
	for i := 0; i < n; i++ {
		id := uuid.New()
		d := fmt.Sprintf("Testing string #%d", i+1)
		kvs[id] = d
		keys[i] = id
	}

	//Add the map to Redis
	if err := cr.SetManyS(r, context.Background(), kvs); err != nil {
		t.Fatal(err)
	}

	//Query the database for the items
	items, err := cr.GetManyS(r, context.Background(), keys...)
	if err != nil {
		t.Fatal(err)
	}

	//Ensure what came in equals what came out
	for i, itm := range items {
		expected := kvs[keys[i]]
		if expected != itm {
			t.Fatalf("item #%d: %v != %v", i+1, expected, itm)
		}
		fmt.Printf("Item #%d: %v\n", i+1, itm)
	}

	//Cleanup what was inserted
	_, err = cr.Del(r, context.Background(), keys...)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRedisKVSet(t *testing.T) {
	//Set the number of items to insert
	n := 5

	//Get a Redis client
	r := redisInit()

	//Create a map of test strings and IDs to go with them
	kvs := make(map[uuid.UUID]Foo)
	keys := make([]uuid.UUID, n)
	for i := 0; i < n; i++ {
		id := mongoutil.MustNewUUID4()
		obj := Foo{id, fmt.Sprintf("Name_%d", i+1), time.Now().Round(0), []string{fmt.Sprintf("ff_%d", (i+1)*2)}}
		kvs[obj.ID.UUID] = obj
		keys[i] = id.UUID
	}

	//Add the map to Redis
	if err := cr.SetMany(r, context.Background(), kvs); err != nil {
		t.Fatal(err)
	}

	//Query the database for the items
	items, err := cr.GetMany[uuid.UUID, Foo](r, context.Background(), keys...)
	if err != nil {
		t.Fatal(err)
	}

	//Ensure what came in equals what came out
	for i, itm := range items {
		expected := kvs[keys[i]]
		if !fooeq(expected, itm) {
			t.Fatalf("item #%d: %v != %v", i+1, expected, itm)
		}
		fmt.Printf("Item #%d: %v\n", i+1, itm)
	}

	//Cleanup what was inserted
	_, err = cr.Del(r, context.Background(), keys...)
	if err != nil {
		t.Fatal(err)
	}
}
