package agp

import (
	"go.mongodb.org/mongo-driver/bson"
)

//
//-- CLASS: AggregationBuilder
//

/*
Represents a MongoDB aggregation pipeline, consisting of multiple stages.
Each stage is a well-formed BSON document (`bson.D`) that defines an
aggregation function and any associated operations or data it needs to
function.
*/
type AggregationBuilder struct {
	stages bson.A
}

// Creates a new blank aggregation pipeline.
func NewAggPipeline() *AggregationBuilder {
	obj := AggregationBuilder{}
	obj.stages = bson.A{}
	return &obj
}

// Returns the full aggregation pipeline as a BSON array (`bson.A`).
func (ab AggregationBuilder) Build() bson.A {
	return ab.stages
}
