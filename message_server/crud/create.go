package crud

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"wraith.me/message_server/db/mongoutil"
	credis "wraith.me/message_server/redis"
)

/*
Creates a database entry in both MongoDB and Redis, using a UUID as the
object ID for MongoDB and the key in Redis. Objects are marshalled to
BSON for MongoDB and GOB for Redis. This function performs a create
operation in CRUD.
*/
func Create[T any](
	//Database drivers & context
	c *mongo.Collection, r *redis.Client, ctx context.Context,
	//Key/value
	id mongoutil.UUID, obj *T,
) (upd UpdatedCount, err error) {
	//Prep: Marshal the incoming object to a BSON document
	bs, err := bson.Marshal(obj)
	if err != nil {
		return UpdatedCount{0, 0}, err
	}

	//Step 1: Insert the document into MongoDB
	_, err = c.InsertOne(ctx, bs)
	if err != nil {
		return UpdatedCount{0, 0}, err
	}

	//Step 2: Cache the keypair in Redis
	if err = credis.Create(r, ctx, id.UUID, obj); err != nil {
		return UpdatedCount{0, 0}, err
	}

	//Step 3: No errors occurred, so return true
	return UpdatedCount{1, 1}, nil
}

//TODO: CreateMany, CreateArray
