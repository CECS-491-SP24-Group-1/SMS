package tests

import (
	"fmt"
	"testing"
	"time"

	ccrypto "wraith.me/message_server/crypto"
	c "wraith.me/message_server/obj/challenge"
	"wraith.me/message_server/util"
)

func TestCTokenClaims(t *testing.T) {
	//Create a new ed25519 key for encryption
	_, sk, _ := ccrypto.NewKeypair(nil)

	//Test the email claim functionality
	etok := c.NewEmailChallenge(
		util.MustNewUUID7(),
		util.MustNewUUID7(),
		c.CPurposeUNKNOWN,
		time.Now().Add(24*time.Hour),
		"jdoe@example.com",
	)
	fmt.Printf("%v\n", etok)
	fmt.Printf("enc: %s\n", etok.Encrypt(sk))
	fmt.Println()

	//Test the pk token claim functionality
	pk, _, err := ccrypto.NewKeypair(nil)
	if err != nil {
		t.Fatal(err)
	}
	pktok := c.NewPKChallenge(
		util.MustNewUUID7(),
		util.MustNewUUID7(),
		c.CPurposeUNKNOWN,
		time.Now().Add(24*time.Hour),
		pk,
	)
	fmt.Printf("%v\n", pktok)
	fmt.Printf("enc: %s\n", pktok.Encrypt(sk))
	fmt.Println()
}
