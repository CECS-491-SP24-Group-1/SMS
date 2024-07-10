package tests

import (
	"fmt"
	"testing"
	"time"

	"wraith.me/message_server/obj/ip_addr"
	"wraith.me/message_server/obj/token"
	"wraith.me/message_server/util"
)

func TestTokenSerialization(t *testing.T) {
	//Create a new token for testing
	tok := token.NewToken(
		util.MustNewUUID7(), ip_addr.FromString("127.0.0.1"),
		token.TokenScopeUSER, time.Now().Add(5*time.Minute),
	)

	//Serialize to a byte array and deserialize it back
	bytes := tok.ToBytes()
	obj := token.TokenFromBytes(bytes)
	//obj.Expire = false //This line causes the test to fail

	//Print the token strings
	fmt.Printf("TOK_IN : %s\n", tok)
	fmt.Printf("TOK_OUT: %s\n", obj)

	//Ensure the tokens are equal
	if !tok.Equal(*obj) {
		t.Fatalf("Tokens are unequal")
	}
}
