package tests

import (
	"crypto/ed25519"
	"testing"

	ccrypto "wraith.me/message_server/crypto"
)

func TestKey(t *testing.T) {
	//Generate a new Ed25519 private key
	edpk, edsk, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatal(err)
	}

	//Convert the private/public keys to our representations
	ccrypto.MustKeyFromBytes(ccrypto.PrivkeyFromBytes, edsk)
	ccrypto.MustKeyFromBytes(ccrypto.PubkeyFromBytes, edpk)

	//Attempt to derive the public key
	//pks
}
