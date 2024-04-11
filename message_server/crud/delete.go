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
Deletes a database entry from both MongoDB and Redis, using a UUID as the
object ID for MongoDB and the key in Redis. This function performs a delete
operation in CRUD.
*/
func Delete(
	//Database drivers & context
	c *mongo.Collection, r *redis.Client, ctx context.Context,
	//Key/value
	id mongoutil.UUID,
) (success bool, err error) {
	//1: Remove the document from MongoDB
	_, err = c.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}

	//2: Remove the keypair from Redis
	_, err = credis.Del(r, ctx, id.UUID)
	if err != nil {
		return false, err
	}

	//3: No errors occurred, so return true
	return true, nil
}

//TODO: DeleteMany, DeleteByFilter
