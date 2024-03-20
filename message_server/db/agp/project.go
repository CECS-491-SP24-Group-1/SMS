package agp

import "go.mongodb.org/mongo-driver/bson"

/*
Defines a `$project` pipeline stage. See the following MongoDB site for
more info: https://www.mongodb.com/docs/manual/reference/operator/aggregation/project/
*/
func (ab *AggregationBuilder) Project(filter bson.D) *AggregationBuilder {
	//Construct the pipeline stage
	stage := bson.D{{Key: "$project", Value: filter}}

	//Append the stage to the pipeline and return the pipeline object for further use
	ab.stages = append(ab.stages, stage)
	return ab
}

/*
Defines a `$project` pipeline stage. This function takes multiple BSON
elements, combines them into a single document, and adds them to the
aggregation pipeline as a single unit. See the following MongoDB site
for more info: https://www.mongodb.com/docs/manual/reference/operator/aggregation/project/
*/
func (ab *AggregationBuilder) ProjectE(filters ...bson.E) *AggregationBuilder {
	//Collect the filters and combine them into a single document
	filtersDoc := bson.D{}
	for _, filter := range filters {
		filtersDoc = append(filtersDoc, filter)
	}

	//Add the stage to the pipeline
	return ab.Project(filtersDoc)
}

// Species a field that should be excluded from a projection.
func P_Hide(key string) bson.E {
	return bson.E{Key: key, Value: 0}
}

// Species a field that should be included in a projection.
func P_Show(key string) bson.E {
	return bson.E{Key: key, Value: 1}
}
