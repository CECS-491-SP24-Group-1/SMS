package db

import "github.com/qiniu/qmgo"

// QMgoBase is the base struct for all collections.
type QMgoBase struct {
	*qmgo.Collection
}

// QMgoCollection interface defines methods that each collection should implement.
type QMgoCollection interface {
	//Gets the parent db name.
	ParentDB() string

	//Gets the name of the collection.
	CollectionName() string
}
