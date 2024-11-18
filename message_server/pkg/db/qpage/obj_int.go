package qpage

import "go.mongodb.org/mongo-driver/bson"

// Defines the structure of a sort param for the filtering process.
type sorter struct {
	Name  string
	Order int
}

// Defines the structure of a returned aggregation op.
type aresult struct {
	Total int64    `bson:"total"`
	Data  []bson.D `bson:"data"`
}
