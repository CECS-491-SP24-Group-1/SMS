package tests

import (
	"fmt"
	"testing"
	"time"

	ccrypto "wraith.me/message_server/crypto"
	"wraith.me/message_server/db/mongoutil"
	c "wraith.me/message_server/obj/challenge"
)

func TestCTokenClaims(t *testing.T) {
	//Test the email claim functionality
	etok := c.NewEmailChallenge(
		mongoutil.MustNewUUID7(),
		mongoutil.MustNewUUID7(),
		c.CPurposeUNKNOWN,
		time.Now().Add(24*time.Hour),
		"jdoe@example.com",
	)
	fmt.Printf("%v\n", etok)

	//Test the pk token claim functionality
	pk, _, err := ccrypto.NewKeypair(nil)
	if err != nil {
		t.Fatal(err)
	}
	pktok := c.NewPKChallenge(
		mongoutil.MustNewUUID7(),
		mongoutil.MustNewUUID7(),
		c.CPurposeUNKNOWN,
		time.Now().Add(24*time.Hour),
		pk,
	)
	fmt.Printf("%v\n", pktok)
}
