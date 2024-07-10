package tests

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"wraith.me/message_server/obj"
	"wraith.me/message_server/obj/ip_addr"
	"wraith.me/message_server/util"
)

var user = obj.NewUserSimple(
	"johndoe123",
	"johndoeismyname@example.com",
)

func TestUser2BSON(t *testing.T) {
	//Add a test token
	tok := obj.NewToken(
		util.MustNewUUID7(), ip_addr.FromString("127.0.0.1"),
		obj.TokenScopeUSER, time.Now().Add(5*time.Minute),
	)
	user.Tokens = append(user.Tokens, *tok)

	//Marshal to BSON
	bb, err := bson.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	//Unmarshal the BSON back to an object
	var out obj.User
	if err := bson.Unmarshal(bb, &out); err != nil {
		t.Fatal(err)
	}

	//Print the input and output objects
	fmt.Printf("IN:  %v\n", *user)
	fmt.Printf("OUT: %v\n", out)
}

func TestUser2JSON(t *testing.T) {
	//Add a test token
	tok := obj.NewToken(
		util.MustNewUUID7(), ip_addr.FromString("127.0.0.1"),
		obj.TokenScopeUSER, time.Now().Add(5*time.Minute),
	)
	user.Tokens = append(user.Tokens, *tok)

	//Marshal to JSON
	jb, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("JSON: %s\n", jb) //Tokens are not marshaled with the JSON

	//Unmarshal the JSON back to an object
	var out obj.User
	if err := json.Unmarshal(jb, &out); err != nil {
		t.Fatal(err)
	}

	//Print the input and output objects
	fmt.Printf("IN:  %v\n", *user)
	fmt.Printf("OUT: %v\n", out)
}
