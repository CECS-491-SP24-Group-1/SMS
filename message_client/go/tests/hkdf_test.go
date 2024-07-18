package tests

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"
	cc "wraith.me/clientside_crypto/crypto"
)

func TestEd25519HKDF(t *testing.T) {
	//Define the passphrase, salt, and ctx
	passphrase := "password12345"
	salt := uuid.MustParse("a91742ba-8771-45f2-92b4-e233d4615438")

	//Set a custom hash function
	//cc.HKDFHashFunc = md5.New
	cc.HKDFHashFunc = sha3.New256

	//Run HKDF
	privkey, err := cc.Ed25519HKDF(passphrase, salt[:], nil)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("pass: %s\n", passphrase)
	fmt.Printf("salt: %s\n", base64.StdEncoding.EncodeToString(salt[:]))
	fmt.Printf("priv: %s\n", privkey.String())
}
