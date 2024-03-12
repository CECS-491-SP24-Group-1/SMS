package tests

import (
	"fmt"
	"net"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"wraith.me/message_server/db/mongoutil"
	"wraith.me/message_server/obj"
)

func TestUser2BSON(t *testing.T) {
	uuid, _ := mongoutil.NewUUID7()
	user := obj.User{
		ID:          *uuid,
		Username:    "john.doe",
		DisplayName: "John Doe",
		Email:       "jdoe.example.com",
		Pubkey:      [obj.PUBKEY_SIZE]byte{},
		LastLogin:   time.Now(),
		LastIP:      net.ParseIP("127.0.0.1"),
		Flags:       obj.DefaultUserFlags(),
	}

	bson, berr := bson.Marshal(user)
	if berr != nil {
		t.Error(berr)
		t.FailNow()
	}
	fmt.Printf("BSON: %s\n", string(bson))
}
