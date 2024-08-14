package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"wraith.me/message_server/config"
	"wraith.me/message_server/db"
	"wraith.me/message_server/obj/token"
	"wraith.me/message_server/schema/user"
	"wraith.me/message_server/util"
)

func TestGenAccessTok(t *testing.T) {
	//Get the server secrets
	secrets, err := config.EnvInit("../secrets.env")
	if err != nil {
		t.Fatal(err)
	}

	//Connect to the database and get the user collection
	if _, err := db.GetInstance().Connect(db.DefaultMConfig()); err != nil {
		t.Fatal(err)
	}
	ucoll := user.GetCollection()

	//Construct a query to get a random user's ID
	query := bson.A{
		bson.D{{Key: "$sample", Value: bson.D{{Key: "size", Value: 1}}}},
		bson.D{{Key: "$project", Value: bson.D{{Key: "_id", Value: 1}}}},
	}

	//Get a random user's ID; the result is deserialized to a map; qmgo doesn't like deserializing to raw objects
	var idres map[string]util.UUID
	err = ucoll.Aggregate(context.Background(), query).One(&idres)
	if err != nil {
		t.Fatal(err)
	}
	id := util.MustGetSingular(idres)

	//Generate an access token
	//The tokens generated should work for auth tests, but are not persisted to the db
	tok := token.NewToken(
		id,
		secrets.ID,
		token.TokenTypeACCESS,
		time.Now().Add(30*time.Minute),
		nil,
	).Encrypt(secrets.SK, true)
	fmt.Printf("User UUID:  %s\n", id)
	fmt.Printf("User Token: %s\n", tok)
}
