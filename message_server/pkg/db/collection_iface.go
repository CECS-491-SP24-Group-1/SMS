package db

// QMgoCollection interface defines methods that each collection should implement.
type QMgoCollection interface {
	//Gets the parent db name.
	ParentDB() string

	//Gets the name of the collection.
	CollectionName() string
}
