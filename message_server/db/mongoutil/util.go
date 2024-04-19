package mongoutil

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Returns an array of matched documents from the database, given an aggregation pipeline.
func Aggregate(coll *mongo.Collection, filter bson.A, ctx context.Context) ([]bson.M, error) {
	var docs []bson.M
	err := AggregateT(&docs, coll, filter, ctx)
	return docs, err

	/*
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
	*/
}

// Returns an array of matched documents from the database, given a string aggregation pipline.
func AggregateS(coll *mongo.Collection, filter string, ctx context.Context) ([]bson.M, error) {
	//Convert the string into a BSON aggregation array
	agg, serr := String2Aggregation(filter)
	if serr != nil {
		return nil, serr
	}

	//Perform the aggregation
	return Aggregate(coll, agg, ctx)
}

/*
Returns an array of matched documents from the database, given an aggregation
pipeline. Additionally, this function outputs the response in a format specified
by `T`.
*/
func AggregateT[T any](ret *T, coll *mongo.Collection, filter bson.A, ctx context.Context) error {
	//Execute the aggregation against the database
	cursor, err := coll.Aggregate(ctx, filter)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	//Collect all the results into a slice
	if err := cursor.All(context.TODO(), ret); err != nil {
		return err
	}
	return nil
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

// Retrieves a list of documents by their IDs.
func FindById[T any](coll *mongo.Collection, ctx context.Context, ids ...UUID) ([]T, error) {
	//Construct the output slice to hold the results
	res := []T{}

	//Construct a BSON document from the list of IDs
	query := bson.M{"_id": bson.M{"$in": ids}}

	//Execute a find operation on the database using the BSON document
	cursor, err := coll.Find(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	//Collect all the results into a slice and return them
	if err := cursor.All(context.TODO(), &res); err != nil {
		return nil, err
	}
	return res, nil
}

/*
Outputs a BSON array, representing a database aggregation pipeline, given
a well-formed ExtJSON string. This function can be chained with `Aggregate()`
to perform database queries from stringified Mongo aggregation piplines.
This simplifies the way queries can be declared, allowing for convenient
and simple edits. This function makes use of `bson.UnmarshalExtJSON()` to
operate. See: https://pkg.go.dev/go.mongodb.org/mongo-driver@latest/bson#UnmarshalExtJSON
*/
func String2Aggregation(str string) (bson.A, error) {
	var agg bson.A
	if err := bson.UnmarshalExtJSON([]byte(str), false, &agg); err != nil {
		return nil, err
	}
	return agg, nil
}
