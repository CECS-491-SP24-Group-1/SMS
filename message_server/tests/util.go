package tests

import (
	"reflect"

	"github.com/qiniu/qmgo"
	"github.com/redis/go-redis/v9"
	"wraith.me/message_server/db"
	cr "wraith.me/message_server/redis"
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
