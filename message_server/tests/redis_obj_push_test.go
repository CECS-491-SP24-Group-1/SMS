package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	cr "wraith.me/message_server/pkg/redis"
)

func TestRedisObjPush(t *testing.T) {
	//Create some instances plus IDs for each
	foo1 := "Hello world"
	foo1ID := uuid.New()
	foo2 := 3.14159
	foo2ID := uuid.New()
	foo3 := 25863168973206
	foo3ID := uuid.New()

	//Connect to Redis
	red := redisInit()

	//Push the objects into the Redis database
	if err := cr.Set(red, context.Background(), foo1ID, &foo1); err != nil {
		t.Fatal(err)
	}
	if err := cr.Set(red, context.Background(), foo2ID, &foo2); err != nil {
		t.Fatal(err)
	}
	if err := cr.Set(red, context.Background(), foo3ID, &foo3); err != nil {
		t.Fatal(err)
	}

	//Get the objects from Redis
	var rfoo1 string
	if err := cr.Get(red, context.Background(), foo1ID, &rfoo1); err != nil {
		t.Fatal(err)
	}
	if foo1 != rfoo1 {
		fmt.Printf("foo1 : %v\n", foo1)
		fmt.Printf("rfoo1: %v\n", rfoo1)
		t.Fatal("foo1 != rfoo1")
	}
	var rfoo2 float64
	if err := cr.Get(red, context.Background(), foo2ID, &rfoo2); err != nil {
		t.Fatal(err)
	}
	if foo2 != rfoo2 {
		fmt.Printf("foo2 : %v\n", foo2)
		fmt.Printf("rfoo2: %v\n", rfoo2)
		t.Fatal("foo2 != rfoo2")
	}
	var rfoo3 int
	if err := cr.Get(red, context.Background(), foo3ID, &rfoo3); err != nil {
		t.Fatal(err)
	}
	if foo3 != rfoo3 {
		fmt.Printf("foo3 : %v\n", foo3)
		fmt.Printf("rfoo3: %v\n", rfoo3)
		t.Fatal("foo3 != rfoo3")
	}
}

func TestRedisCObjPush(t *testing.T) {
	//Connect to Redis
	red := redisInit()

	//Push the objects into the Redis database
	if err := cr.Set(red, context.Background(), foo1.ID.UUID, &foo1); err != nil {
		t.Fatal(err)
	}
	if err := cr.Set(red, context.Background(), foo2.ID.UUID, &foo2); err != nil {
		t.Fatal(err)
	}
	if err := cr.Set(red, context.Background(), foo3.ID.UUID, &foo3); err != nil {
		t.Fatal(err)
	}

	//Get the objects from Redis
	var rfoo1 Foo
	if err := cr.Get(red, context.Background(), foo1.ID.UUID, &rfoo1); err != nil {
		t.Fatal(err)
	}
	if !fooeq(foo1, rfoo1) {
		fmt.Printf("foo1 : %v\n", foo1)
		fmt.Printf("rfoo1: %v\n", rfoo1)
		t.Fatal("foo1 != rfoo1")
	}
	var rfoo2 Foo
	if err := cr.Get(red, context.Background(), foo2.ID.UUID, &rfoo2); err != nil {
		t.Fatal(err)
	}
	if !fooeq(foo2, rfoo2) {
		fmt.Printf("foo2 : %v\n", foo2)
		fmt.Printf("rfoo2: %v\n", rfoo2)
		t.Fatal("foo2 != rfoo2")
	}
	var rfoo3 Foo
	if err := cr.Get(red, context.Background(), foo3.ID.UUID, &rfoo3); err != nil {
		t.Fatal(err)
	}
	if !fooeq(foo3, rfoo3) {
		fmt.Printf("foo3 : %v\n", foo3)
		fmt.Printf("rfoo3: %v\n", rfoo3)
		t.Fatal("foo3 != rfoo3")
	}
}

func TestRedisCObjDel(t *testing.T) {
	//Connect to Redis
	red := redisInit()

	//Push the objects into the Redis database
	if err := cr.Set(red, context.Background(), foo1.ID.UUID, &foo1); err != nil {
		t.Fatal(err)
	}

	//Get the objects from Redis
	var rfoo1 Foo
	if err := cr.Get(red, context.Background(), foo1.ID.UUID, &rfoo1); err != nil {
		t.Fatal(err)
	}
	if !fooeq(foo1, rfoo1) {
		fmt.Printf("foo1 : %v\n", foo1)
		fmt.Printf("rfoo1: %v\n", rfoo1)
		t.Fatal("foo1 != rfoo1")
	}

	//Delete the object and ensure no errors occurred
	if _, err := cr.Del(red, context.Background(), foo1.ID.UUID); err != nil {
		t.Fatal(err)
	}
}

func TestRedisCObjMod(t *testing.T) {
	//Connect to Redis
	red := redisInit()

	//Push the objects into the Redis database
	if err := cr.Set(red, context.Background(), foo1.ID.UUID, &foo1); err != nil {
		t.Fatal(err)
	}

	//Get the objects from Redis
	var rfoo1 Foo
	if err := cr.Get(red, context.Background(), foo1.ID.UUID, &rfoo1); err != nil {
		t.Fatal(err)
	}
	if !fooeq(foo1, rfoo1) {
		fmt.Printf("foo1 : %v\n", foo1)
		fmt.Printf("rfoo1: %v\n", rfoo1)
		t.Fatalf("foo1 != rfoo1")
	}

	//Create an updated object
	foo2 := foo1
	foo2.Name = "Jane Doe"
	foo2.FavoriteFoods = []string{"pancakes", "waffles"}

	fmt.Printf("before: %v\n", foo1)
	fmt.Printf("after:  %v\n", foo2)

	//Push the updated object to the database
	if err := cr.Set(red, context.Background(), foo2.ID.UUID, &foo2); err != nil {
		t.Fatal(err)
	}

	//Ensure the change went through successfully
	var rfoo2 Foo
	if err := cr.Get(red, context.Background(), foo1.ID.UUID, &rfoo2); err != nil {
		t.Fatal(err)
	}
	if !fooeq(foo2, rfoo2) {
		fmt.Printf("foo2 : %v\n", foo2)
		fmt.Printf("rfoo2: %v\n", rfoo2)
		t.Fatal("foo2 != rfoo2")
	}
}
