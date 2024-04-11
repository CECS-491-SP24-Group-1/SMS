package tests

import (
	"fmt"
	"testing"
	"time"

	"wraith.me/message_server/db/mongoutil"
	"wraith.me/message_server/obj"
	"wraith.me/message_server/obj/ip_addr"
)

func TestTokenSerialization(t *testing.T) {
	//Create a new token for testing
	tok := obj.NewToken(
		mongoutil.MustNewUUID7(), ip_addr.FromString("127.0.0.1"),
		obj.TokenScopeUSER, time.Now().Add(5*time.Minute),
	)

	//Serialize to a byte array and deserialize it back
	bytes := tok.ToBytes()
	obj := obj.TokenFromBytes(bytes)
	//obj.Expire = false //This line causes the test to fail

	//Print the token strings
	fmt.Printf("TOK_IN : %s\n", tok)
	fmt.Printf("TOK_OUT: %s\n", obj)

	//Ensure the tokens are equal
	if !tok.Equal(*obj) {
		t.Errorf("Tokens are unequal")
		t.FailNow()
	}
}
