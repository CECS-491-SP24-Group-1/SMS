package user

import (
	"sync"

	"wraith.me/message_server/pkg/db"
)

var (
	// Holds the shared instance of this collection.
	userCollectionInst *UserCollection

	// Guard mutex to ensure that only one singleton object is created.
	userCollectionOnce sync.Once
)

/*
Represents a single `User` object in a collection of objects in the database.
This collection is managed by the `qmgo` Mongo ODM library.
*/
type UserCollection struct {
	*db.QMgoBase
}

// This line enforces UserCollection to implement db.QMgoCollection.
var _ db.QMgoCollection = (*UserCollection)(nil)

func (uc UserCollection) ParentDB() string {
	return db.ROOT_DB
}

func (uc UserCollection) CollectionName() string {
	return db.USERS_COLLECTION
}

/*
Gets the currently active collection object instance or initializes it.
This can be safely called multiple times in the program to ensure a
non-nil instance of the collection due to the usage of `sync.Once` to
initialize the singleton.
*/
func GetCollection() *UserCollection {
	userCollectionOnce.Do(func() {
		c := db.GetCollectionManager().GetCollection(UserCollection{})
		userCollectionInst = &UserCollection{c}
	})
	return userCollectionInst
}
