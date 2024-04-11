package tests

import (
	"context"
	"testing"

	"wraith.me/message_server/crud"
)

func TestCrudCreateOne(t *testing.T) {
	//Connect to Mongo and Redis
	m := mongoInit()
	r := redisInit()

	//Get the target collection
	coll := m.Database(TEST_DB_NAME).Collection(TEST_COLL_NAME)

	//Create the document
	_, err := crud.Create[Foo](
		coll, r, context.Background(),
		foo1.ID, &foo1,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCrudDeleteOne(t *testing.T) {
	//Connect to Mongo and Redis
	m := mongoInit()
	r := redisInit()

	//Get the target collection
	coll := m.Database(TEST_DB_NAME).Collection(TEST_COLL_NAME)

	//Create the document
	_, err := crud.Create[Foo](
		coll, r, context.Background(),
		foo1.ID, &foo1,
	)
	if err != nil {
		t.Fatal(err)
	}

	//Delete the document
	_, err = crud.Delete(
		coll, r, context.Background(),
		foo1.ID,
	)
	if err != nil {
		t.Fatal(err)
	}
}
