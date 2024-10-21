package chatroom

import (
	"sync"

	"wraith.me/message_server/pkg/db"
)

var (
	// Holds the shared instance of this collection.
	roomCollectionInst *RoomCollection

	// Guard mutex to ensure that only one singleton object is created.
	roomCollectionOnce sync.Once
)

/*
Represents a single `Room` object in a collection of objects in the database.
This collection is managed by the `qmgo` Mongo ODM library.
*/
type RoomCollection struct {
	*db.QMgoBase
}

// This line enforces RoomCollection to implement db.QMgoCollection.
var _ db.QMgoCollection = (*RoomCollection)(nil)

func (uc RoomCollection) ParentDB() string {
	return db.ROOT_DB
}

func (uc RoomCollection) CollectionName() string {
	return db.CROOMS_COLLECTION
}

/*
Gets the currently active collection object instance or initializes it.
This can be safely called multiple times in the program to ensure a
non-nil instance of the collection due to the usage of `sync.Once` to
initialize the singleton.
*/
func GetCollection() *RoomCollection {
	roomCollectionOnce.Do(func() {
		c := db.GetCollectionManager().GetCollection(RoomCollection{})
		roomCollectionInst = &RoomCollection{c}
	})
	return roomCollectionInst
}
