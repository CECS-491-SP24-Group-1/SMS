package crud

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"wraith.me/message_server/db"
	"wraith.me/message_server/db/mongoutil"
)

/*
Retrieves a list of a user's tokens from either Redis if cached or from
MongoDB if Redis doesn't have them. If a cache miss occurs, then this
function will automatically cache them in Redis for faster future
retrievals. The tokens are encoded as a string, and must be parsed to be
usable. This operation is equivalent to an R operation in CRUD.
*/
func GetSTokens(
	//Database drivers & context
	m *mongo.Client, r *redis.Client, ctx context.Context,
	//ID of the target user
	uid mongoutil.UUID,
) ([]string, error) {
	//Step 1: Check if Redis has tokens for the user
	rt, err := r.LRange(ctx, uid.String(), 0, -1).Result()
	if err != nil && err != redis.Nil {
		//An error occurred with Redis that isn't a null-key error; bail out
		fmt.Printf("[TOK_CRUD_R; REDIS ERROR]: %s\n", err)
		return nil, err
	} else if len(rt) > 0 && err != redis.Nil {
		//Cache hit! Return the tokens from Redis
		//fmt.Printf("CACHE HIT - err: %s, len: %d\n", "<nil>", len(rt))
		return rt, nil
	}

	/*
		At this point, it is safe to assume that a cache miss occurred. Redis either
		returned a `redis.Nil` error or returned an array with no tokens. MongoDB
		should be queried for the token list instead.
	*/
	//fmt.Printf("CACHE MISS - err: `%s`, len: %d\n", err, len(rt))

	//Step 2a: Define the Mongo query and object to load the tokens into
	query := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: uid}}}},
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "_id", Value: 0},
					{Key: "tokens", Value: 1},
				},
			},
		},
	}
	type tokenRet struct {
		Tokens []string `bson:"tokens"`
	}

	//Step 2b: Execute the MongoDB query. Should return a singular `tokenRet` in the output array if the database has tokens for the subject
	ucoll := m.Database(db.ROOT_DB).Collection(db.USERS_COLLECTION)
	var mr []tokenRet
	if err := mongoutil.AggregateT(&mr, ucoll, query, ctx); err != nil {
		fmt.Printf("[TOK_CRUD_R; MONGO AGGREGATION ERROR]: %s\n", err)
		return nil, err
	}

	//Step 3: Check if the MongoDB query returned nothing; no need to cache if there isn't anything to cache
	if len(mr) == 0 || len(mr[0].Tokens) == 0 {
		return mr[0].Tokens, nil
	}

	//Step 4: Cache the tokens in Redis; further requests with the same user ID will be honored by Redis instead
	if err := r.RPush(ctx, uid.String(), mr[0].Tokens).Err(); err != nil {
		//Report the Redis cache error but do not fail; caching not working isn't too big of a deal; just silently fail
		fmt.Printf("[TOK_CRUD_R; REDIS CACHE RPUSH ERROR]: %s\n", err)
	}

	//Step 5: Return the MongoDB tokens to the user, along with a nil error
	return mr[0].Tokens, nil
}

/*
TODO: Add the below functions:
func GetSTokenById; R
func AddSTokens; C
func AddSToken; C
func RevokeSTokens; D
func RevokeSTokenById; D
PLUS non-string versions (maybe not)
* Tokens are immutable and do not have a U operation
*/
