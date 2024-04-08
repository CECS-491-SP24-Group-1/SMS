package tests

import (
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"wraith.me/message_server/db"
	credis "wraith.me/message_server/redis"
)

func mongoInit() *mongo.Client {
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
	rcfg := credis.DefaultRConfig()
	rclient, rerr := credis.GetInstance().Connect(rcfg)
	if rerr != nil {
		panic(rerr)
	}
	return rclient
}
