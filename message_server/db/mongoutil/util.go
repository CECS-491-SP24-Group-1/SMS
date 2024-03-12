package mongoutil

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Returns an array of matched documents from the database, given an aggregation pipline.
func Aggregate(coll *mongo.Collection, filter bson.A, ctx context.Context) ([]bson.M, error) {
	//Execute the aggregation against the database
	cursor, err := coll.Aggregate(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	//Collect all the results into a slice
	var results []bson.M
	if err := cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	//Return the slice
	return results, nil
}

// Checks if a field exists in the database.
func Exists(coll *mongo.Collection, filter bson.D, ctx context.Context) bool {
	//Run the query on the database
	hit := coll.FindOne(ctx, filter)
	var bson bson.M

	//Attempt to decode the document; this will throw an error if the query didn't return a document
	//Thus, return true when no error occurs and false when there is one
	return hit.Decode(&bson) == nil
}
