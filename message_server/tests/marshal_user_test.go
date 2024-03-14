package tests

import (
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"wraith.me/message_server/obj"
)

var user, _ = obj.NewUserSimple(
	"johndoe123",
	"johndoeismyname@example.com",
)

func TestUser2BSON(t *testing.T) {
	//Marshal to BSON
	bb, err := bson.Marshal(user)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	//Unmarshal the BSON back to an object
	var out obj.User
	if err := bson.Unmarshal(bb, &out); err != nil {
		t.Error(err)
		t.FailNow()
	}

	//Print the input and output objects
	fmt.Printf("IN:  %v\n", *user)
	fmt.Printf("OUT: %v\n", out)
}
