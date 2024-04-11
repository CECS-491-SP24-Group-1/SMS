package crud

/*
Contains a pairing of two integers, representing the number of documents
updated in Mongo and Redis.
*/
type UpdatedCount struct {
	//The number of documents updated in MongoDB.
	M int64

	//The number of documents updated in Redis.
	R int64
}
