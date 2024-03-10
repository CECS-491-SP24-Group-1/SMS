package lib

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestEd25519Keygen(t *testing.T) {
	//Generate a keypair
	kp := Ed25519Keygen()

	//Print stuff
	fmt.Printf("SK: %s\n", hex.EncodeToString(kp.SK[:]))
	fmt.Printf("PK: %s\n", hex.EncodeToString(kp.PK[:]))
	fmt.Printf("FP: %s\n", kp.Fingerprint)

	//JSON
	json := kp.JSON()
	obj, _ := Ed25519FromJSON(json)
	fmt.Printf("JSON: %s\n", json)
	fmt.Printf("Obj: %s\n", obj)
}
