package tests

import (
	"context"
	"fmt"
	"log"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"wraith.me/message_server/pkg/db"
	"wraith.me/message_server/pkg/util"
)

const (
	MDB_HOST = "127.0.0.1"
	MDB_PORT = 27017
)

type Test struct {
	//ID  string `bson:"_id"`
	ID  util.UUID `bson:"_id"`
	Foo int       `bson:"foo"`
	Bar string    `bson:"bar"`
	Baz struct {
		Bash float64 `bson:"bash"`
		Ash  [5]int  `bson:"ash"`
	} `bson:"baz"`
}

func DefaultTest() Test {
	obj := Test{}

	obj.ID = util.MustNewUUID7()
	obj.Foo = 500
	obj.Bar = "hello world"
	obj.Baz.Bash = 3.14
	obj.Baz.Ash = [5]int{5, 4, 3, 2, 1}

	return obj
}

func TestInit(t *testing.T) {
	//Connect to MongoDB
	mcfg := db.DefaultMConfig()
	//mcfg.Username = ""
	//mcfg.Password = ""
	client, err := db.GetInstance().Connect(mcfg)
	if err != nil {
		panic(err)
	}
	defer db.GetInstance().Disconnect()

	/*
		//Test listing
		names, _ := client.ListDatabaseNames(context.TODO(), bson.M{})
		for i, name := range names {
			fmt.Printf("DB #%d: %s\n", i, name)
		}
	*/

	//Create sample data
	testobj := DefaultTest()
	sample1, err := bson.Marshal(testobj)
	if err != nil {
		log.Panic(err)
	}

	//Push the data to the database
	collection := client.Database(db.ROOT_DB).Collection(db.TESTS_COLLECTION)
	insertResult, err := collection.InsertOne(context.Background(), sample1)
	if err != nil {
		panic(err)
	}
	fmt.Println(insertResult)

	id := testobj.ID.String()
	fmt.Printf("Object pushed has UUID: %s\n", id)
}
