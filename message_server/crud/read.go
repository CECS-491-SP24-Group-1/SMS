package crud

/*
Performs a CRUD read operation involving MongoDB and Redis. The basic flow
of control is as follows:

- check if the key exists in Redis, henceforth denoted as a cache hit/miss

- if there is a cache hit, serve the object using Redis

- if there's a cache miss, do the following

--
*/
/*
func ReadArray[T any](
	m *mongo.Client, r *redis.Client,
	ctx context.Context,
	key mongoutil.UUID,

) (res []T, err error) {
	//Setup the variables needed to hold the results

	//Query the database and get the tokens of the subject; Redis is used here to cache the results
	redisTokens, err := r.LRange(ctx, key.String(), 0, -1).Result()
	if err == redis.Nil || len(redisTokens) == 0 {
		//Cache miss; query MongoDB for the tokens of the subject and add them to Redis for later
		//Define the aggregation to run
		query := bson.A{
			bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: key}}}},
			bson.D{
				{Key: "$project",
					Value: bson.D{
						{Key: "_id", Value: 0},
						{Key: "tokens", Value: 1},
					},
				},
			},
		}

		//Define the structure of the output; a list of string tokens
		type tokenRet struct {
			Tokens []string `bson:"tokens"`
		}

		//Execute the aggregation. Should return a singular `tokenRet` in the output array if the database has tokens for the subject
		ucoll := m.Database(db.ROOT_DB).Collection(db.USERS_COLLECTION)
		var results []tokenRet
		if aerr := mongoutil.AggregateT(&results, ucoll, query, ctx); aerr != nil {
			//If there is a database aggregation error, deny the client by default for safety; do not report the error
			httpu.HttpErrorAsJson(w, fmt.Errorf("auth; %s", ErrAuthGeneric), http.StatusUnauthorized)
			fmt.Printf("[AUTH; MONGO AGGREGATION ERROR]: %s\n", aerr)
			return
		}

		//If the query returns nothing, simply bail out; the subject has no tokens to begin with
		if len(results) == 0 || len(results[0].Tokens) == 0 {
			httpu.HttpErrorAsJson(w, fmt.Errorf("auth; %s", ErrAuthNoTokenFound), http.StatusUnauthorized)
			return
		}

		//Copy the tokens to the subject tokens array
		res = results[0].Tokens

		/*
			Add the subject ID and corresponding tokens array to Redis. Further
			requests with the same token subject will be honored by Redis instead
			of MongoDB (cache hit). If any errors occur, silently fail and do not
			alert the client. Authentication can still be done, but caching won't
			work.
		* /
		cacheSetResult := amw.rclient.RPush(r.Context(), tokSubject.String(), results[0].Tokens)
		if err := cacheSetResult.Err(); err != nil {
			fmt.Printf("[AUTH; REDIS CACHE RPUSH ERROR]: %s\n", err)
		}
	} else if err != nil {
		//If there is a cache retrieval error, deny the client by default for safety; do not report the error
		httpu.HttpErrorAsJson(w, fmt.Errorf("auth; %s", ErrAuthGeneric), http.StatusUnauthorized)
		fmt.Printf("[AUTH; REDIS ERROR]: %s\n", err)
		return
	} else {
		//Cache hit; copy the token list from Redis
		subjectTokens = redisTokens
	}
}
*/
