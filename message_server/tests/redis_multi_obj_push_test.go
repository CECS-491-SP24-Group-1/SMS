package tests

import (
	"context"
	"fmt"
	"slices"
	"testing"

	"github.com/google/uuid"
	cr "wraith.me/message_server/redis"
)

/*
Sets a key and array of values in the Redis database. If the key already
exists, its value is updated. If not, then a new key is created. Applicable
to C, U in CRUD. See: https://stackoverflow.com/a/53697645
*/
/*
func UpdateOne[T any](c *redis.Client, ctx context.Context, key uuid.UUID, value T, idx int) error {
}
*/

func TestRedisMultiObjPush(t *testing.T) {
	//Create some instances plus IDs for each
	foo1 := "Hello world"
	foo2 := "how are you doing"
	foo3 := "nice to meet you"

	//Construct the array and an ID for it
	aid := uuid.New()
	objs := []string{foo1, foo2, foo3}

	//Connect to Redis
	red := redisInit()

	//Push the object list into the Redis database
	if err := cr.CreateSA(red, context.Background(), aid, objs); err != nil {
		t.Fatal(err)
	}

	//Get the objects from Redis
	robjs, err := cr.GetSA(red, context.Background(), aid)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(objs, robjs) {
		fmt.Printf("objs : %v\n", objs)
		fmt.Printf("robjs: %v\n", robjs)
		t.Fatal("objs != robjs")
	}
}

func TestRedisMultiCObjPush(t *testing.T) {
	//Construct the array and an ID for it
	aid := uuid.New()
	objs := []Foo{foo1, foo2, foo3}

	//Connect to Redis
	red := redisInit()

	//Push the object list into the Redis database
	if err := cr.CreateA(red, context.Background(), aid, objs); err != nil {
		t.Fatal(err)
	}
	//Get the objects from Redis
	robjs, err := cr.GetA[Foo](red, context.Background(), aid)
	if err != nil {
		t.Fatal(err)
	}
	if !fooeqa(objs, robjs) {
		fmt.Printf("objs : %v\n", objs)
		fmt.Printf("robjs: %v\n", robjs)
		t.Fatal("objs != robjs")
	}
}

func TestRedisMultiCObjDel(t *testing.T) {
	//Construct the array and an ID for it
	aid := uuid.New()
	objs := []Foo{foo1}

	//Connect to Redis
	red := redisInit()

	//Push the object list into the Redis database
	if err := cr.CreateA(red, context.Background(), aid, objs); err != nil {
		t.Fatal(err)
	}
	//Get the objects from Redis
	robjs, err := cr.GetA[Foo](red, context.Background(), aid)
	if err != nil {
		t.Fatal(err)
	}
	if !fooeqa(objs, robjs) {
		fmt.Printf("objs : %v\n", objs)
		fmt.Printf("robjs: %v\n", robjs)
		t.Fatal("objs != robjs")
	}

	//Delete the object and ensure no errors occurred
	if _, err := cr.Del(red, context.Background(), aid); err != nil {
		t.Fatal(err)
	}
}

func TestRedisMultiCObjModOne(t *testing.T) {
	//Construct the array and an ID for it
	aid := uuid.New()
	objs := []Foo{foo1, foo2, foo3}

	//Connect to Redis
	red := redisInit()

	//Push the object list into the Redis database
	if err := cr.CreateA(red, context.Background(), aid, objs); err != nil {
		t.Fatal(err)
	}
	//Get the objects from Redis
	robjs, err := cr.GetA[Foo](red, context.Background(), aid)
	if err != nil {
		t.Fatal(err)
	}
	if !fooeqa(objs, robjs) {
		fmt.Printf("objs : %v\n", objs)
		fmt.Printf("robjs: %v\n", robjs)
		t.Fatal("objs != robjs")
	}

	//Create an updated object
	foo3a := foo3
	foo3a.Name = "James Doe"
	foo3a.FavoriteFoods = []string{"tic-tacs", "ice-cream"}
	midx := int64(2)

	fmt.Printf("before: %v\n", foo3)
	fmt.Printf("after:  %v\n", foo3a)

	//Push the updated object to the database
	if err := cr.SetAt(red, context.Background(), aid, midx, foo3a); err != nil {
		t.Fatal(err)
	}

	//Ensure the change went through successfully
	var rfoo3a Foo
	if err := cr.GetAt(red, context.Background(), aid, midx, &rfoo3a); err != nil {
		t.Fatal(err)
	}
	if !fooeq(foo3a, rfoo3a) {
		fmt.Printf("foo3a : %v\n", foo3a)
		fmt.Printf("rfoo3a: %v\n", rfoo3a)
		t.Fatal("foo3a != rfoo3a")
	}
}
