package lib

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

const (
	ED25519_LEN = 32 //Ed25519 keys are always 32 bytes long.
)

// Represents an Ed25519 keypair.
type Ed25519KP struct {
	SK          [ED25519_LEN]byte `json:"sk"`          //Holds the private key.
	PK          [ED25519_LEN]byte `json:"pk"`          //Holds the public key.
	Fingerprint string            `json:"fingerprint"` //Holds the fingerprint of the public key as a SHA-256 hash.
}

// Generates a new Ed25519 keypair.
func Ed25519Keygen() Ed25519KP {
	pubkey, privkey, _ := ed25519.GenerateKey(nil)
	return Ed25519FromBytes(privkey.Seed(), pubkey)
}

// Derives an Ed25519 keypair object from raw bytes.
func Ed25519FromBytes(sk []byte, pk []byte) Ed25519KP {
	//Create the base object
	out := Ed25519KP{}

	//Assign the public and private key bytes
	copy(out.SK[:], sk[:])
	copy(out.PK[:], pk[:])

	//Hash the public key
	h := sha256.Sum256(pk)
	out.Fingerprint = hex.EncodeToString(h[:])

	//Return the object
	return out
}

// Derives an Ed25519 keypair object from a JSON string.
func Ed25519FromJSON(jsons string) (obj Ed25519KP, err error) {
	//Attempt to unmarshal an object from the JSON
	err = json.Unmarshal([]byte(jsons), &obj)
	return
}

// Derives an Ed25519 keypair object from a private key via `scalar_mult()“.
func Ed25519FromSK(sk []byte) Ed25519KP {
	//Get the public key equivalent via `scalar_mult()`
	pubSmult := ed25519.NewKeyFromSeed(sk).Public()

	//Return the object
	return Ed25519FromBytes(sk, []byte(pubSmult.(ed25519.PublicKey)))
}

// Checks if this Ed25519 keypair is equal to another.
func (kp Ed25519KP) Equal(other Ed25519KP) bool {
	return bytes.Equal(kp.SK[:], other.SK[:]) && bytes.Equal(kp.PK[:], other.PK[:])
}

// Returns the JSON representation of the object.
func (kp Ed25519KP) JSON() string {
	json, _ := json.Marshal(kp)
	return string(json)
}

// Returns the string representation of the object.
func (kp Ed25519KP) String() string {
	return fmt.Sprintf("Ed25519KP{sk=%s, pk=%s}", hex.EncodeToString(kp.SK[:]), hex.EncodeToString(kp.PK[:]))
}
