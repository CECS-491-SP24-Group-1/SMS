package tests

import (
	"encoding/json"
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

func TestUser2JSON(t *testing.T) {
	//Marshal to JSON
	jb, err := json.Marshal(user)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	//Unmarshal the JSON back to an object
	var out obj.User
	if err := json.Unmarshal(jb, &out); err != nil {
		t.Error(err)
		t.FailNow()
	}

	//Print the input and output objects
	fmt.Printf("IN:  %v\n", *user)
	fmt.Printf("OUT: %v\n", out)
}

func TestRand(t *testing.T) {
	type Test struct {
		Status obj.ChallengeStatus
	}
	status := Test{Status: obj.ChallengeStatusPENDING}

	jb, err := bson.Marshal(status)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Printf("Status: `%v`\n", jb)
}
