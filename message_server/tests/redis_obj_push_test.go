package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	credis "wraith.me/message_server/redis"
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
	if err := credis.Set(red, context.Background(), foo1ID, &foo1); err != nil {
		t.Error(err)
	}
	if err := credis.Set(red, context.Background(), foo2ID, &foo2); err != nil {
		t.Error(err)
	}
	if err := credis.Set(red, context.Background(), foo3ID, &foo3); err != nil {
		t.Error(err)
	}

	//Get the objects from Redis
	var rfoo1 string
	if err := credis.Get(red, context.Background(), foo1ID, &rfoo1); err != nil {
		t.Error(err)
	}
	if foo1 != rfoo1 {
		fmt.Printf("foo1 : %v\n", foo1)
		fmt.Printf("rfoo1: %v\n", rfoo1)
		t.Errorf("foo1 != rfoo1")
	}
	var rfoo2 float64
	if err := credis.Get(red, context.Background(), foo2ID, &rfoo2); err != nil {
		t.Error(err)
	}
	if foo2 != rfoo2 {
		fmt.Printf("foo2 : %v\n", foo2)
		fmt.Printf("rfoo2: %v\n", rfoo2)
		t.Errorf("foo2 != rfoo2")
	}
	var rfoo3 int
	if err := credis.Get(red, context.Background(), foo3ID, &rfoo3); err != nil {
		t.Error(err)
	}
	if foo3 != rfoo3 {
		fmt.Printf("foo3 : %v\n", foo3)
		fmt.Printf("rfoo3: %v\n", rfoo3)
		t.Errorf("foo3 != rfoo3")
	}
}

func TestRedisCObjPush(t *testing.T) {
	//Connect to Redis
	red := redisInit()

	//Push the objects into the Redis database
	if err := credis.Set(red, context.Background(), foo1.ID, &foo1); err != nil {
		t.Error(err)
	}
	if err := credis.Set(red, context.Background(), foo2.ID, &foo2); err != nil {
		t.Error(err)
	}
	if err := credis.Set(red, context.Background(), foo3.ID, &foo3); err != nil {
		t.Error(err)
	}

	//Get the objects from Redis
	var rfoo1 Foo
	if err := credis.Get(red, context.Background(), foo1.ID, &rfoo1); err != nil {
		t.Error(err)
	}
	if !fooeq(foo1, rfoo1) {
		fmt.Printf("foo1 : %v\n", foo1)
		fmt.Printf("rfoo1: %v\n", rfoo1)
		t.Errorf("foo1 != rfoo1")
	}
	var rfoo2 Foo
	if err := credis.Get(red, context.Background(), foo2.ID, &rfoo2); err != nil {
		t.Error(err)
	}
	if !fooeq(foo2, rfoo2) {
		fmt.Printf("foo2 : %v\n", foo2)
		fmt.Printf("rfoo2: %v\n", rfoo2)
		t.Errorf("foo2 != rfoo2")
	}
	var rfoo3 Foo
	if err := credis.Get(red, context.Background(), foo3.ID, &rfoo3); err != nil {
		t.Error(err)
	}
	if !fooeq(foo3, rfoo3) {
		fmt.Printf("foo3 : %v\n", foo3)
		fmt.Printf("rfoo3: %v\n", rfoo3)
		t.Errorf("foo3 != rfoo3")
	}
}

func TestRedisCObjDel(t *testing.T) {
	//Connect to Redis
	red := redisInit()

	//Push the objects into the Redis database
	if err := credis.Set(red, context.Background(), foo1.ID, &foo1); err != nil {
		t.Error(err)
	}

	//Get the objects from Redis
	var rfoo1 Foo
	if err := credis.Get(red, context.Background(), foo1.ID, &rfoo1); err != nil {
		t.Error(err)
	}
	if !fooeq(foo1, rfoo1) {
		fmt.Printf("foo1 : %v\n", foo1)
		fmt.Printf("rfoo1: %v\n", rfoo1)
		t.Errorf("foo1 != rfoo1")
	}

	//Delete the object and ensure no errors occurred
	if _, err := credis.Del(red, context.Background(), foo1.ID); err != nil {
		t.Error(err)
	}
}

func TestRedisCObjMod(t *testing.T) {
	//Connect to Redis
	red := redisInit()

	//Push the objects into the Redis database
	if err := credis.Set(red, context.Background(), foo1.ID, &foo1); err != nil {
		t.Error(err)
	}

	//Get the objects from Redis
	var rfoo1 Foo
	if err := credis.Get(red, context.Background(), foo1.ID, &rfoo1); err != nil {
		t.Error(err)
	}
	if !fooeq(foo1, rfoo1) {
		fmt.Printf("foo1 : %v\n", foo1)
		fmt.Printf("rfoo1: %v\n", rfoo1)
		t.Errorf("foo1 != rfoo1")
	}

	//Create an updated object
	foo2 := foo1
	foo2.Name = "Jane Doe"
	foo2.FavoriteFoods = []string{"pancakes", "waffles"}

	fmt.Printf("before: %v\n", foo1)
	fmt.Printf("after:  %v\n", foo2)

	//Push the updated object to the database
	if err := credis.Set(red, context.Background(), foo2.ID, &foo2); err != nil {
		t.Error(err)
	}

	//Ensure the change went through successfully
	var rfoo2 Foo
	if err := credis.Get(red, context.Background(), foo1.ID, &rfoo2); err != nil {
		t.Error(err)
	}
	if !fooeq(foo2, rfoo2) {
		fmt.Printf("foo2 : %v\n", foo2)
		fmt.Printf("rfoo2: %v\n", rfoo2)
		t.Errorf("foo2 != rfoo2")
	}
}
