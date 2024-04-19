package crud

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"wraith.me/message_server/db"
	"wraith.me/message_server/db/mongoutil"
	chall "wraith.me/message_server/obj/challenge"
	cr "wraith.me/message_server/redis"
	"wraith.me/message_server/util"
)

/*
Adds a single challenge or list of challenges to the database. This function
also automatically caches the new entries in Redis, using the Redis CRUD
wrappers, which marshals them to a binary string. The number of challenges
added is returned as an integer. This operation is equivalent to a C operation
in CRUD.
*/
func AddChallenges(
	//Database drivers & context
	m *mongo.Client, r *redis.Client, ctx context.Context,
	//Challenge list
	chl ...*chall.Challenge,
) (int, error) {
	//Step 1: Create arrays to hold the Mongo write queue and Redis byte strings
	bwq := make([]mongo.WriteModel, len(chl))
	kvs := make(map[uuid.UUID]string)

	//Step 2: Loop over the list of challenges to insert
	for i, ch := range chl {
		//Step 2a: Marshal the current challenge to a byte string and add it to the kv map
		bs, err := util.ToGOBBytes(ch)
		if err != nil {
			return -i, err
		}
		kvs[ch.ID.UUID] = string(bs)

		//Step 2b: Marshal the current challenge to BSON
		doc, err := bson.Marshal(ch)
		if err != nil {
			return -i, err
		}

		//Step 2c: Define the Mongo filters for targeting and add the document to the bulk write array
		//See: https://stackoverflow.com/a/66489583
		mfilter := bson.D{{Key: "_id", Value: ch.ID}}
		bwq[i] = mongo.NewReplaceOneModel().SetFilter(mfilter).SetReplacement(doc).SetUpsert(true)
	}

	//Step 3: Insert the documents into MongoDB via bulk write
	coll := m.Database(db.ROOT_DB).Collection(db.CHALL_COLLECTION)
	res, err := coll.BulkWrite(ctx, bwq)
	if err != nil {
		return defualtErrAmt, err
	}

	//Step 4: Cache the challenges in Redis
	if err := cr.SetManyS(r, ctx, kvs); err != nil {
		return defualtErrAmt, err
	}

	//Step 5a: Pick the modified quantity for Mongo, picking the biggest between inserted, updated, and upserted
	mcount := max(res.InsertedCount, res.ModifiedCount, res.UpsertedCount)

	//Step 5b: No errors occurred, so return the number of documents inserted
	return int(mcount), nil
}

/*
Removes a single challenge or list of challenges from the database. This
function also removes the entries from Redis, using the Redis CRUD wrappers.
The number of challenges removed is returned as an integer. This operation
is equivalent to a D operation in CRUD.
*/
func RemoveChallenges(
	//Database drivers & context
	m *mongo.Client, r *redis.Client, ctx context.Context,
	//Challenge ID list
	ids ...mongoutil.UUID,
) (int, error) {
	//Step 1a: Define the Mongo query that targets all challenges specified in the varargs by ID
	targets := bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: ids}}}}

	//Step 1b: Remove the documents from MongoDB
	coll := m.Database(db.ROOT_DB).Collection(db.CHALL_COLLECTION)
	res, err := coll.DeleteMany(ctx, targets)
	if err != nil {
		return defualtErrAmt, err
	}

	//Step 2: Remove the tokens from Redis
	//TODO: use the redis crud wrapper here
	sids := make([]string, len(ids))
	for i, v := range ids {
		sids[i] = v.String()
	}
	if err := r.Del(ctx, sids...).Err(); err != nil {
		return defualtErrAmt, err
	}

	//Step 3: No errors occurred, so return the number of documents deleted
	return int(res.DeletedCount), nil
}

func GetChallengesById(
	//Database drivers & context
	m *mongo.Client, r *redis.Client, ctx context.Context,
	//Challenge ID list
	ids ...mongoutil.UUID,
) ([]chall.Challenge, error) {
	//Step 1: Check if Redis has the challenges
	ch, rerr := cr.GetMany[mongoutil.UUID, chall.Challenge](r, ctx, ids...)
	if rerr.Cause() != nil && rerr.Cause() != redis.Nil {
		//An error occurred with Redis that isn't a null-key error; bail out
		fmt.Printf("[CHALL_CRUD_R; REDIS ERROR]: %s\n", rerr)
		return nil, rerr
	} else if len(ch) > 0 && rerr.Cause() != redis.Nil {
		//Cache hit! Return the challenges from Redis
		fmt.Printf("CACHE HIT - err: %s, len: %d\n", "<nil>", len(ch))
		return ch, nil
	}

	/*
		At this point, it is safe to assume that a cache miss occurred with at least
		one queried challenge. Redis returned a `redis.Nil` error or an array with no
		challenges. MongoDB should be queried for any challenges that weren't found,
		based on the indices in the "problematic indices" array.
	*/
	fmt.Printf("CACHE MISS - err: `%s`, len: %d\n", rerr, len(ch))

	//Step 2: Get the IDs of the challenges that weren't in Redis
	missedIds := []mongoutil.UUID{}
	for _, mid := range rerr.Indices() {
		missedIds = append(missedIds, ids[mid])
	}

	//Step 3: Submit the query to MongoDB and get the missing challenges
	coll := m.Database(db.ROOT_DB).Collection(db.CHALL_COLLECTION)
	challs, err := mongoutil.FindById[chall.Challenge](coll, ctx, missedIds...)
	if err != nil {
		return nil, err
	}

	//Step 4: Check if the MongoDB query returned nothing; no need to cache if there isn't anything to cache
	if len(challs) == 0 {
		return challs, nil
	}

	//Step 5a: Create a mapping of the missing challenges to their ID
	cmap := make(map[mongoutil.UUID]chall.Challenge)
	for i, chall := range challs {
		cmap[missedIds[i]] = chall
	}

	fmt.Printf("chs: %+v\n", cmap)

	//Step 5b: Cache the missing challenges in Redis; further requests with the same challenge IDs will be honored by Redis instead

	//Step 1: Create an output array for the challenges that's the same size as the input list
	//out := make([]chall.Challenge, len(ids))

	/*
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
	*/
	return nil, nil
}

/*
TODO: Add the below functions:
func GetChallenges; R
func GetChallengeById; R
func AddChallenges; C
func AddChallenge; C
func RemoveChallenges; D
func RemoveChallengeById; D
* Challenges are immutable and do not have a U operation
*/
