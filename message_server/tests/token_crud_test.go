package tests

import (
	"context"
	"fmt"
	"testing"

	"wraith.me/message_server/crud"
	"wraith.me/message_server/db/mongoutil"
)

func TestTokCRUDGet(t *testing.T) {
	//Get a Mongo and Redis client
	m := mongoInit()
	r := redisInit()

	//Query for tokens by the user's ID
	uid := mongoutil.UUIDFromStringOrNil("018eac6a-9a9d-77f2-bd4b-238b34d9c767")
	toks, err := crud.GetSTokens(m, r, context.Background(), uid)
	if err != nil {
		t.Fatal(err)
	}

	//Print the array of tokens
	fmt.Printf("Tokens for user with UUID %s:\n", uid)
	fmt.Printf("%v", toks)
}
