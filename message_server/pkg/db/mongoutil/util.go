package mongoutil

import (
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"wraith.me/message_server/pkg/util"
)

// Converts an aggregation result of `_id: UUID` to an array.
func Aggregate2IDArr(agg qmgo.AggregateI) ([]util.UUID, error) {
	//Perform the aggregation
	mp := make([]map[string]util.UUID, 0) //Array of single valued maps, keyed by `_id`
	err := agg.All(&mp)
	if err != nil {
		return nil, err
	}

	//Loop over the aggregation array and collect all UUIDs to an array of just UUIDs
	out := make([]util.UUID, len(mp))
	for i, uuid := range mp {
		out[i] = uuid["_id"]
	}
	return out, nil
}

// Creates a BSON array out of a slice.
func Slice2BsonA[T any](slice []T) bson.A {
	bsonArray := make(bson.A, len(slice))
	for i, s := range slice {
		bsonArray[i] = s
	}
	return bsonArray
}
