package tests

import (
	"context"
	"reflect"

	"github.com/qiniu/qmgo"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"wraith.me/message_server/pkg/db"
	cr "wraith.me/message_server/pkg/redis"
	"wraith.me/message_server/pkg/schema/user"
)

const (
	TEST_DB_NAME   = "tests"
	TEST_COLL_NAME = "tc"
)

// Determines if an incoming object is a complex type, ie: a struct.
func isComplexType[T any](targ T) bool {
	//Get the type and kind of the target
	t := reflect.TypeOf(targ)
	k := t.Kind()

	//For arrays and slices, set the kind to be the first item, if available
	if (k == reflect.Array || k == reflect.Slice) && reflect.ValueOf(targ).Len() > 0 {
		k = reflect.ValueOf(targ).Index(0).Kind()
	}

	//Only return true for structs and interfaces
	return k == reflect.Struct || k == reflect.Interface
}

func mongoInit() *qmgo.Client {
	//Connect to MongoDB
	mcfg := db.DefaultMConfig()
	mclient, merr := db.GetInstance().Connect(mcfg)
	if merr != nil {
		panic(merr)
	}
	return mclient
}

func redisInit() *redis.Client {
	//Connect to Redis
	rcfg := cr.DefaultRConfig()
	rclient, rerr := cr.GetInstance().Connect(rcfg)
	if rerr != nil {
		panic(rerr)
	}
	return rclient
}

// Gets a random user from the database.
func GetRandomUser() (*user.User, error) {
	//Get the database instance
	dbc := db.GetInstance()

	//Connect if not already done so
	if !dbc.IsConnected() {
		if _, err := dbc.Connect(db.DefaultMConfig()); err != nil {
			return nil, err
		}
	}
	ucoll := user.GetCollection()

	//Construct a query to get a random user
	query := bson.A{
		bson.D{{Key: "$sample", Value: bson.D{{Key: "size", Value: 1}}}},
	}

	//Get a random user; the result is deserialized to a map; qmgo doesn't like deserializing to raw objects
	var res user.User
	err := ucoll.Aggregate(context.Background(), query).One(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
