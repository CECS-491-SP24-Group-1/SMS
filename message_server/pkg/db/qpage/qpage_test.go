package qpage

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"testing"
	"time"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/maps"
	"wraith.me/message_server/pkg/db/mongoutil"
	"wraith.me/message_server/pkg/util"
)

type todo struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Desc      string             `json:"desc" bson:"desc"`
	Owner     string             `json:"owner" bson:"owner"`
	Completed bool               `json:"completed" bson:"completed"`
}

// Name::IsMale
var names = map[string]bool{
	"Furina":      false,
	"Zhongli":     true,
	"Xilonen":     false,
	"Neuvillette": true,
	"Kazuha":      true,
	"Bennett":     true,
	"Navia":       false,
	"Arlecchino":  false,
	"Chiori":      false,
	"Yelan":       false,
}

func NewTodo(desc, owner string) todo {
	//Generate a new object ID
	oid := primitive.NewObjectID()

	//Generate values for blank fields
	if desc == "" {
		desc = fmt.Sprintf("Untitled TODO with ID %s", oid.Hex())
	}
	if owner == "" {
		//Get a random number
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(names))))
		if err != nil {
			panic(fmt.Sprintf("NewTodo::pickAWinner: %s", err))
		}

		//Pick a winner
		owner = maps.Keys(names)[n.Int64()]
	}

	//Create the object
	return todo{
		ID:        oid,
		Desc:      desc,
		Owner:     owner,
		Completed: false,
	}
}

func TestAddTodo_x1(t *testing.T) {
	//Setup the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coll := connectDB(ctx)

	//Add the todo
	todo := NewTodo("", "")
	_, err := coll.InsertOne(ctx, todo)
	if err != nil {
		t.Fatalf("failed to insert TODO item %s: %s", todo.ID.Hex(), err)
	}
}

func TestAddTodo_x10(t *testing.T) {
	//Setup the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coll := connectDB(ctx)

	//Add 10 todos
	for i := 0; i < 10; i++ {
		//Add the todo
		todo := NewTodo("", "")
		_, err := coll.InsertOne(ctx, todo)
		if err != nil {
			t.Fatalf("failed to insert TODO item #%d %s: %s", i+1, todo.ID.Hex(), err)
		}
	}
}

func TestPurgeTodos(t *testing.T) {
	// Setup the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coll := connectDB(ctx)

	//Drop the collection
	if err := coll.DropCollection(ctx); err != nil {
		t.Fatalf("failed to drop todos collection: %s", err)
	}
}

func TestAggregation(t *testing.T) {
	//Setup the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coll := connectDB(ctx)

	//Setup the aggregation pipeline
	pipeline := bson.A{
		bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "owner", Value: bson.D{
					{Key: "$concat", Value: bson.A{
						bson.D{{Key: "$toString", Value: "$owner"}},
						"_test",
					}},
				}},
			}},
		},
	}

	//Setup the paginator
	paginate, err := NewQPage(coll)
	if err != nil {
		t.Fatal(err)
	}

	//Set the pagination params
	//skipId, _ := primitive.ObjectIDFromHex("6736b0abbdc1c6abbfd313df")
	params := Params{
		Page:    6,
		PerPage: 75,
		//SkipToID:     skipId,
	}

	//Perform the query
	todos := make([]todo, 0)
	pagination, err := paginate.Aggregate(&todos, ctx, pipeline, params)
	if err != nil {
		log.Fatal(err)
	}

	for i, doc := range todos {
		fmt.Printf("Doc #%d: %v\n", i, doc)
	}
	fmt.Printf("Total: %+v\n", pagination)
}

func TestAggregation2NewType(t *testing.T) {
	//Define the custom type to unmarshal to
	type customType struct {
		ID     primitive.ObjectID `json:"id" bson:"_id"`
		Owner  string             `json:"owner" bson:"owner"`
		Gender string             `json:"gender" bson:"gender"`
	}

	//Setup the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coll := connectDB(ctx)

	//Get the list of male and female characters from the input map
	males := mongoutil.Slice2BsonA(
		util.GetKeysByValue(names, true),
	)
	females := mongoutil.Slice2BsonA(
		util.GetKeysByValue(names, false),
	)

	//Setup the aggregation pipeline
	pipeline := bson.A{
		bson.D{{"$project", bson.D{{"desc", 0}, {"completed", 0}}}},
		bson.D{{"$set", bson.D{{"gender", bson.D{{"$switch", bson.D{{"branches", bson.A{
			bson.D{{"case", bson.D{{"$in", bson.A{"$owner", males}}}}, {"then", "Male"}},
			bson.D{{"case", bson.D{{"$in", bson.A{"$owner", females}}}}, {"then", "Female"}},
		}}, {"default", "UNKNOWN"}}}}}}}},
	}

	//Setup the paginator
	paginate, err := NewQPage(coll)
	if err != nil {
		t.Fatal(err)
	}

	//Set the pagination params
	//skipId, _ := primitive.ObjectIDFromHex("6736b0abbdc1c6abbfd313df")
	params := Params{
		Page:    6,
		PerPage: 75,
		//SkipToID:     skipId,
	}

	//Perform the query
	todos := make([]customType, 0)
	pagination, err := paginate.Aggregate(&todos, ctx, pipeline, params)
	if err != nil {
		log.Fatal(err)
	}

	for i, doc := range todos {
		fmt.Printf("Doc #%d: %v\n", i, doc)
	}
	fmt.Printf("Total: %+v\n", pagination)
}

func TestFind(t *testing.T) {
	//Setup the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coll := connectDB(ctx)

	//Setup the find query
	query := bson.D{{Key: "owner", Value: "Furina"}}

	//Setup the paginator
	paginate, err := NewQPage(coll)
	if err != nil {
		t.Fatal(err)
	}

	//Set the pagination params
	//skipId, _ := primitive.ObjectIDFromHex("6736b0abbdc1c6abbfd313df")
	params := Params{
		Page:    3,
		PerPage: 75,
		//SkipToID:     skipId,
	}

	//Perform the query
	todos := make([]todo, 0)
	pagination, err := paginate.Find(&todos, ctx, query, params)
	if err != nil {
		log.Fatal(err)
	}

	for i, doc := range todos {
		fmt.Printf("Doc #%d: %v\n", i, doc)
	}
	fmt.Printf("Total: %+v\n", pagination)

}

func connectDB(ctx context.Context) *qmgo.Collection {
	//mgoCfg := &qmgo.Config{Uri: "mongodb://localhost:27017", Database: "testdb", Coll: "todos"}
	mgoCfg := &qmgo.Config{Uri: "mongodb://localhost:27017"}
	client, err := qmgo.NewClient(ctx, mgoCfg)
	if err != nil {
		log.Fatal(err)
	}
	return client.Database("testdb").Collection("todos")
}
