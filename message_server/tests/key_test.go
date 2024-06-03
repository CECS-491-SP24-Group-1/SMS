package tests

import (
	"crypto/ed25519"
	"fmt"
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

	sk, _ := ccrypto.ParsePrivkeyBytes("eK7Rv8dfHPrWgeVcHIoskqMNke2EjWUFaIgafCaU3ZE=")
	pk, _ := ccrypto.ParsePubkeyBytes("wDw04q6c94g7zn5IwGe1M0E6NJRDuHCa0x+joia8DFg=")
	fmt.Printf("sk: `%v`\n", sk)

	msg := "this is a test"
	sig := ccrypto.Sign(sk, []byte(msg))
	fmt.Printf("sig: '%v'\n", sig)
	fmt.Printf("ok: '%v'\n", ccrypto.Verify(pk, []byte(msg), sig))
}
