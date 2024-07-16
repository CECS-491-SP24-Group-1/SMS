package tests

import (
	"context"
	"fmt"
	"testing"

	"wraith.me/message_server/db"
	"wraith.me/message_server/db/agp"
)

func TestStringifiedMongoAgg(t *testing.T) {
	//Create the aggregation string
	str := `[{"$match":{"$or":[{"username":"angelina59"},{"email":"timothy.volkman10@yahoo.com"}]}},{"$project":{"_id":1}}]`

	//Connect to the database
	mcfg := db.DefaultMConfig()
	mclient, merr := db.GetInstance().Connect(mcfg)
	if merr != nil {
		t.Fatal(merr)
	}
	defer db.GetInstance().Disconnect()

	//Run the query
	ucoll := mclient.Database(db.ROOT_DB).Collection(db.USERS_COLLECTION)
	res := ucoll.Aggregate(context.Background(), str)
	//res, err := mongoutil.AggregateS(ucoll, str, context.Background())

	fmt.Printf("res: %v\n", res)

	/*
		//Show results
		for i := 0; i < len(res); i++ {
			fmt.Printf("RES #%d: %+v\n", i, res[i])
		}
	*/
}

func TestAggregations(t *testing.T) {
	/*
		//Create test aggregation
		agg := bson.A{
			bson.D{
				{Key: "$match",
					Value: bson.D{
						{Key: "flags.should_purge", Value: true},
						{Key: "flags.purge_by", Value: bson.D{{Key: "$gt", Value: time.Now()}}},
					},
				},
			},
			bson.D{{Key: "$project", Value: bson.D{{Key: "purge_by", Value: "$flags.purge_by"}}}},
		}
	*/

	agg := agp.NewAggPipeline().
		ProjectE(
			agp.P_Show("flags.should_purge"),
			agp.P_Hide("_id"),
		).Build()

	//Connect to the database
	mcfg := db.DefaultMConfig()
	cli, merr := db.GetInstance().Connect(mcfg)
	if merr != nil {
		t.Fatal(merr)
	}
	defer db.GetInstance().Disconnect()

	//Define the target collection
	coll := cli.Database(db.ROOT_DB).Collection(db.TESTS_COLLECTION)

	//Run the aggregation
	res := coll.Aggregate(context.Background(), agg)

	//Show results
	fmt.Printf("RES: %v\n", res)
}
