package crud

var (
	//The number of documents affected if there was an error.
	errorCount = CRUDCount{
		M: -1, R: -1,
	}
)

/*
Contains a pairing of the number of documents affected by a CRUD operation
for both MongoDB and Redis.
*/
type CRUDCount struct {
	M int64 //The number of documents affected in MongoDB.
	R int64 //The number of documents affected in Redis.
}
