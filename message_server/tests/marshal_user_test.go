package tests

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"wraith.me/message_server/obj/ip_addr"
	"wraith.me/message_server/obj/token"
	"wraith.me/message_server/schema/user"
	"wraith.me/message_server/util"
)

var usr = user.NewUserSimple(
	"johndoe123",
	"johndoeismyname@example.com",
)

func TestUser2BSON(t *testing.T) {
	//Add a test token
	tok := token.NewToken(
		util.MustNewUUID7(), ip_addr.FromString("127.0.0.1"),
		token.TokenScopeUSER, time.Now().Add(5*time.Minute),
	)
	usr.Tokens = append(usr.Tokens, *tok)

	//Marshal to BSON
	bb, err := bson.Marshal(usr)
	if err != nil {
		t.Fatal(err)
	}

	//Unmarshal the BSON back to an object
	var out user.User
	if err := bson.Unmarshal(bb, &out); err != nil {
		t.Fatal(err)
	}

	//Print the input and output objects
	fmt.Printf("IN:  %v\n", *usr)
	fmt.Printf("OUT: %v\n", out)
}

func TestUser2JSON(t *testing.T) {
	//Add a test token
	tok := token.NewToken(
		util.MustNewUUID7(), ip_addr.FromString("127.0.0.1"),
		token.TokenScopeUSER, time.Now().Add(5*time.Minute),
	)
	usr.Tokens = append(usr.Tokens, *tok)

	//Marshal to JSON
	jb, err := json.Marshal(usr)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("JSON: %s\n", jb) //Tokens are not marshaled with the JSON

	//Unmarshal the JSON back to an object
	var out user.User
	if err := json.Unmarshal(jb, &out); err != nil {
		t.Fatal(err)
	}

	//Print the input and output objects
	fmt.Printf("IN:  %v\n", *usr)
	fmt.Printf("OUT: %v\n", out)
}

func TestUser2MS(t *testing.T) {
	//Create a mock user
	user := user.NewUserSimple("johndoe", "johndoe@example.com")
	fmt.Printf("struct: %v\n", user)

	//Marshal to a map using mapstructure
	ms := make(map[string]interface{})
	if err := util.MSTextMarshal(*user, &ms, "bson"); err != nil {
		t.Fatal(err)
	}

	//Redact some fields as a test
	ms["id"] = ms["UUID"]
	delete(ms, "UUID")
	delete(ms, "flags")
	delete(ms, "last_ip")
	delete(ms, "options")
	delete(ms, "tokens")

	//Emit the map to stdout
	fmt.Printf("mapstr: %v\n", ms)
}
