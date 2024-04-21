package crud

import (
	"context"

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
	chl ...chall.Challenge,
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
Gets a series of challenges by their IDs from either Redis if cached or
from MongoDB if Redis doesn't have them. If a cache miss occurs, then this
function will automatically cache them in Redis for faster future retrievals.
A cache hit only occurs if all challenges requested are available. Otherwise,
MongoDB will be consulted for the list of "problematic objects", or those
that weren't found in Redis. This operation is equivalent to an R operation in
CRUD.
*/
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
		return nil, rerr
	} else if len(ch) > 0 && rerr.Cause() != redis.Nil {
		//Cache hit! Return the challenges from Redis
		//fmt.Printf("FULL CACHE HIT - err: %s, len: %d\n", "<nil>", len(ch))
		return ch, nil
	}

	/*
		At this point, it is safe to assume that a cache miss occurred with at least
		one queried challenge. Redis returned a `redis.Nil` error or an array with no
		challenges. MongoDB should be queried for any challenges that weren't found,
		based on the indices in the "problematic indices" array.
	*/
	//fmt.Printf("1+ CACHE MISSES - err: `%s`, len: %d\n", rerr, len(ch))

	//Step 2: Get the IDs of the challenges that weren't in Redis
	missedIds := []mongoutil.UUID{}
	missedIdxs := []int{}
	for _, mid := range rerr.Indices() {
		missedIdxs = append(missedIdxs, mid)
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

	//Step 5a: Create a mapping of the missing challenges to their ID and add them to the output array
	cmap := make(map[mongoutil.UUID]chall.Challenge, len(missedIds))
	for i, mid := range missedIds {
		ch[missedIdxs[i]] = challs[i]
		cmap[mid] = challs[i]
	}

	//Step 5b: Add the missing challenges to Redis
	if err := cr.SetMany(r, ctx, cmap); err != nil {
		return nil, err
	}

	//Step 6: Return the list of challenges
	return ch, nil
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
	_, err = cr.Del(r, ctx, ids...)
	if err != nil {
		return defualtErrAmt, err
	}

	//Step 3: No errors occurred, so return the number of documents deleted
	return int(res.DeletedCount), nil
}
