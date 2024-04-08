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
	credis "wraith.me/message_server/redis"
)

/*
Sets a key and array of values in the Redis database. If the key already
exists, its value is updated. If not, then a new key is created. Applicable
to C, U in CRUD. See: https://stackoverflow.com/a/53697645
*/
func rstructSetA[T any](c *redis.Client, ctx context.Context, key string, values []T) error {
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
	if err := rstructSetA(red, context.Background(), aid.String(), objs); err != nil {
		t.Error(err)
	}

	//Get the objects from Redis
	robjs, err := credis.GetA[string](red, context.Background(), aid.String())
	if err != nil {
		t.Error(err)
	}
	if !slices.Equal(objs, robjs) {
		fmt.Printf("objs : %v\n", objs)
		fmt.Printf("robjs: %v\n", robjs)
		t.Errorf("objs != robjs")
	}
}

func TestRedisMultiCObjPush(t *testing.T) {
	//Define the object
	type Foo struct {
		ID            uuid.UUID
		Name          string
		Birthday      time.Time
		FavoriteFoods []string
	}

	eq := func(a Foo, b Foo) bool {
		return a.ID == b.ID && a.Name == b.Name && a.Birthday == b.Birthday && slices.Equal(a.FavoriteFoods, b.FavoriteFoods)
	}
	eqa := func(a []Foo, b []Foo) bool {
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
	objs := []Foo{foo1, foo2, foo3}

	//Connect to Redis
	red := redisInit()

	//Push the object list into the Redis database
	if err := rstructSetA(red, context.Background(), aid.String(), objs); err != nil {
		t.Error(err)
	}
	//Get the objects from Redis
	robjs, err := credis.GetA[Foo](red, context.Background(), aid.String())
	if err != nil {
		t.Error(err)
	}
	if !eqa(objs, robjs) {
		fmt.Printf("objs : %v\n", objs)
		fmt.Printf("robjs: %v\n", robjs)
		t.Errorf("objs != robjs")
	}
}

func TestRedisMultiCObjDel(t *testing.T) {
	//Define the object
	type Foo struct {
		ID            uuid.UUID
		Name          string
		Birthday      time.Time
		FavoriteFoods []string
	}

	eq := func(a Foo, b Foo) bool {
		return a.ID == b.ID && a.Name == b.Name && a.Birthday == b.Birthday && slices.Equal(a.FavoriteFoods, b.FavoriteFoods)
	}
	eqa := func(a []Foo, b []Foo) bool {
		return slices.EqualFunc(a, b, eq)
	}

	//Create some instances
	foo1 := Foo{
		ID:            uuid.New(),
		Name:          "John Doe",
		Birthday:      time.Now().Round(0),
		FavoriteFoods: []string{"carrots", "apples", "pasta"},
	}

	//Construct the array and an ID for it
	aid := uuid.New()
	objs := []Foo{foo1}

	//Connect to Redis
	red := redisInit()

	//Push the object list into the Redis database
	if err := rstructSetA(red, context.Background(), aid.String(), objs); err != nil {
		t.Error(err)
	}
	//Get the objects from Redis
	robjs, err := credis.GetA[Foo](red, context.Background(), aid.String())
	if err != nil {
		t.Error(err)
	}
	if !eqa(objs, robjs) {
		fmt.Printf("objs : %v\n", objs)
		fmt.Printf("robjs: %v\n", robjs)
		t.Errorf("objs != robjs")
	}

	/*
		//Delete the object and ensure no errors occurred
		if _, err := credis.Del(red, context.Background(), aid.String()); err != nil {
			t.Error(err)
		}
	*/
}
