package user

import (
	"sync"

	"github.com/qiniu/qmgo"
	"wraith.me/message_server/db"
)

const (
	//The name of the database.
	UserCollectionDBName = db.ROOT_DB

	//The name of the collection in the database.
	UserCollectionCollName = db.USERS_COLLECTION
)

/*
Represents a single `User` object in a collection of objects in the database.
This collection is managed by the `qmgo` Mongo ODM library.
*/
type UserCollection struct {
	*qmgo.Collection
}

// Holds the shared instance of this collection.
var _UserCollectionInst *UserCollection

// Guard mutex to ensure that only one singleton object is created.
var _UserCollectionOnce sync.Once

/*
Gets the currently active collection object instance or initializes it.
This can be safely called multiple times in the program to ensure a
non-nil instance of the collection due to the usage of `sync.Once` to
initialize the singleton.
*/
func GetCollection() *UserCollection {
	_UserCollectionOnce.Do(func() {
		//Get the active database client instance
		client := db.GetInstance().GetClient()

		//Set the collection options
		coll := client.Database(UserCollectionDBName).Collection(UserCollectionCollName)

		//Assign the singleton
		_UserCollectionInst = &UserCollection{coll}
	})
	return _UserCollectionInst
}
