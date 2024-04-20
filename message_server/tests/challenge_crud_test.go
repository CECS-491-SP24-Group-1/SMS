package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"wraith.me/message_server/crud"
	"wraith.me/message_server/db/mongoutil"
	"wraith.me/message_server/obj"
	c "wraith.me/message_server/obj/challenge"
)

/*
func TestChallCRUDGet(t *testing.T) {
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
*/

func TestRChallCRUDAddDel(t *testing.T) {
	//Get a Mongo and Redis client
	m := mongoInit()
	r := redisInit()

	//Create some challenges to add
	c1 := c.NewChallenge(
		c.ChallengeScopeEMAIL,
		obj.Identifiable{ID: mongoutil.MustNewUUID4(), Type: obj.IdTypeUSER},
		obj.Identifiable{ID: mongoutil.MustNewUUID4(), Type: obj.IdTypeUSER},
		time.Now().Add(c.DEFAULT_CHALLENGE_EXPIRY),
	)
	c2 := c.NewChallenge(
		c.ChallengeScopeEMAIL,
		obj.Identifiable{ID: mongoutil.MustNewUUID4(), Type: obj.IdTypeUSER},
		obj.Identifiable{ID: mongoutil.MustNewUUID4(), Type: obj.IdTypeUSER},
		time.Now().Add(c.DEFAULT_CHALLENGE_EXPIRY),
	)
	c3 := c.NewChallenge(
		c.ChallengeScopeEMAIL,
		obj.Identifiable{ID: mongoutil.MustNewUUID4(), Type: obj.IdTypeUSER},
		obj.Identifiable{ID: mongoutil.MustNewUUID4(), Type: obj.IdTypeUSER},
		time.Now().Add(c.DEFAULT_CHALLENGE_EXPIRY),
	)

	//Add the challenges to the database via the C CRUD operation
	_, err := crud.AddChallenges(m, r, context.Background(), c1, c2, c3)
	if err != nil {
		t.Fatal(err)
	}

	//List the IDs of the added challenges
	fmt.Println("Added the following challenges to the database")
	fmt.Printf("    - %s; payload: %s\n", c1.ID.String(), c1.Payload)
	fmt.Printf("    - %s; payload: %s\n", c2.ID.String(), c2.Payload)
	fmt.Printf("    - %s; payload: %s\n", c3.ID.String(), c3.Payload)

	//Remove the challenges from the database via the D CRUD operation
	_, err = crud.RemoveChallenges(m, r, context.Background(), c1.ID, c2.ID, c3.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestChallCRUDAdd(t *testing.T) {
	//Get a Mongo and Redis client
	m := mongoInit()
	r := redisInit()

	//Add the challenges to the database via the C CRUD operation
	cnt, err := crud.AddChallenges(m, r, context.Background(), chall1, chall2, chall3)
	if err != nil {
		t.Fatal(err)
	}

	//List the IDs of the added challenges
	fmt.Println("Added the following challenges to the database")
	fmt.Printf("    - %s; payload: %s\n", chall1.ID.String(), chall1.Payload)
	fmt.Printf("    - %s; payload: %s\n", chall2.ID.String(), chall2.Payload)
	fmt.Printf("    - %s; payload: %s\n", chall3.ID.String(), chall3.Payload)
	fmt.Printf("Total: %d\n", cnt)
}

func TestChallCRUDDel(t *testing.T) {
	//Get a Mongo and Redis client
	m := mongoInit()
	r := redisInit()

	//Remove the challenges from the database via the D CRUD operation
	cnt, err := crud.RemoveChallenges(m, r, context.Background(), chall1.ID, chall2.ID, chall3.ID)
	if err != nil {
		t.Fatal(err)
	}

	//List the IDs of the removed challenges
	fmt.Println("Removed the following challenges from the database")
	fmt.Printf("    - %s\n", chall1.ID.String())
	fmt.Printf("    - %s\n", chall2.ID.String())
	fmt.Printf("    - %s\n", chall3.ID.String())
	fmt.Printf("Total: %d\n", cnt)
}

// TODO: test cache durability (if the received challenges still have the same payload for both hits, deletion then get with a miss, then get with a hit)
func TestChallCRUDGet(t *testing.T) {
	//Get a Mongo and Redis client
	m := mongoInit()
	r := redisInit()

	//Attempt to get 2/3 of the challenges by ID from the database
	chs, err := crud.GetChallengesById(m, r, context.Background(), chall1.ID, chall2.ID)
	if err != nil {
		t.Fatal(err)
	}

	//Display the challenges gotten
	fmt.Println("Received the following challenges from the database")
	for _, ch := range chs {
		fmt.Printf("    - %s; payload: %s\n", ch.ID.String(), ch.Payload)
	}
	fmt.Printf("Total: %d\n", len(chs))
}
