package tests

import (
	"crypto/ed25519"
	"fmt"
	"testing"

	ccrypto "wraith.me/message_server/pkg/crypto"
)

func TestKey(t *testing.T) {
	//Generate a new Ed25519 private key
	edpk, edsk, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatal(err)
	}

	//Convert the private/public keys to our representations
	ccrypto.MustFromBytes(ccrypto.PrivkeyFromBytes, edsk)
	ccrypto.MustFromBytes(ccrypto.PubkeyFromBytes, edpk)

	//Attempt to derive the public key
	//pks

	sk, _ := ccrypto.ParsePrivkey("eK7Rv8dfHPrWgeVcHIoskqMNke2EjWUFaIgafCaU3ZE=")
	pk, _ := ccrypto.ParsePubkey("wDw04q6c94g7zn5IwGe1M0E6NJRDuHCa0x+joia8DFg=")
	fmt.Printf("sk: `%v`\n", sk)

	msg := "v4.local.MnyKlOkL6aegNJFuy1YnjEWNlj8EZPj4HCClqkgs5dSqt_WBElKdPz8K55KtqZzcF8yJDU-33LdSH7_oamlK0-GvaYHwe_rs9dvue3aeF_5i9uxMuT9KEMUT5-H_dR-EqoXk8Q71eK-doxETQnlCnB2D7skyAzunN5p9Mtl-daq9gcXoJcl1eBRwCWaqOlOgO2wWscRiIOhq7Zec7fcdT4BmaE2C48BM3nxWyYeX7lwcZl-2VMqWpZWJ3nCMry3mFYzdb663uIm8D7G3HKz8lB5N184UWNw6GrEmEaX_Ob2g_cGQkZFM5O7aPhhB3jgoUurmkJCDaxUgaGc7ZzI1w_08NwmYu5-Z1pmyppc8IZdwfN7CC9T87GRSVlmsqjvxDCax7h6oOrbIFNwI0NRQmsnDeO_FAfLDTOnoYlw4HNhYDR_y6AIHbpxPpMXLwme8j_MZL2TJFOUUbz0FqKNsmESxfuWx3y1PjG9AdqvW_bonTL5AQfU_blY7QHQcvpDKgJkEu4DaFwP6"
	sig := ccrypto.Sign(sk, []byte(msg))
	fmt.Printf("sig: '%v'\n", sig)
	fmt.Printf("ok: '%v'\n", ccrypto.Verify(pk, []byte(msg), sig))
}
