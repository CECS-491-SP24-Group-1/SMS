package db

import (
	"sync"
)

// CollectionManager manages the creation and retrieval of collection instances.
type CollectionManager struct {
	collections map[string]*QMgoBase
	mu          sync.Mutex
}

var (
	manager     *CollectionManager
	managerOnce sync.Once
)

// GetCollectionManager returns the singleton instance of CollectionManager.
func GetCollectionManager() *CollectionManager {
	managerOnce.Do(func() {
		manager = &CollectionManager{
			collections: make(map[string]*QMgoBase),
		}
	})
	return manager
}

// GetCollection returns the collection instance for the given QMgoCollection.
func (cm *CollectionManager) GetCollection(c QMgoCollection) *QMgoBase {
	//Lock the mutex
	cm.mu.Lock()
	defer cm.mu.Unlock()

	//Get the parent db and collection names
	pdb := c.ParentDB()
	cname := c.CollectionName()

	//Get the collection
	//After this point, it is assumed that the collection doesn't exist
	key := pdb + "." + cname
	if coll, exists := cm.collections[key]; exists {
		return coll
	}

	//Get the active database client instance
	client := GetInstance().GetClient()

	//Set the collection options
	coll := client.Database(pdb).Collection(cname)

	//Create and store the new collection instance
	newColl := &QMgoBase{coll}
	cm.collections[key] = newColl

	return newColl
}
