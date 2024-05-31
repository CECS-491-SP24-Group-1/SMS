package crypto

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

const (
	PUBKEY_SIZE = ed25519.PublicKeySize
)

//
//-- ALIAS: Pubkey
//

// Represents the bytes of an entity's public key.
type Pubkey [PUBKEY_SIZE]byte

// Gets the fingerprint of a `Pubkey` object using SHA256.
func (pkb Pubkey) Fingerprint() string {
	hash := sha256.Sum256(pkb[:])
	return hex.EncodeToString(hash[:])
}

// Marshals a `Pubkey` object to JSON.
func (pkb Pubkey) MarshalJSON() ([]byte, error) {
	return json.Marshal(pkb.String())
}

// Parses a `Pubkey` object from a string.
func ParsePubkeyBytes(str string) (*Pubkey, error) {
	//Derive a byte array from the string
	ba, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}

	//Ensure the byte array length is correct
	if len(ba) != PUBKEY_SIZE {
		return nil, fmt.Errorf("mismatched byte array size (%d); expected: %d", len(ba), PUBKEY_SIZE)
	}

	//Copy the bytes to a new object and return it
	obj := &Pubkey{}
	copy(obj[:], ba)
	return obj, nil
}

// Converts a `Pubkey` object to a string.
func (pkb Pubkey) String() string {
	return base64.StdEncoding.EncodeToString(pkb[:])
}

// Unmarshals a `Pubkey` object from JSON.
func (pkb *Pubkey) UnmarshalJSON(b []byte) error {
	//Unmarshal to a string
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	//Derive a valid object from the string and reassign
	obj, err := ParsePubkeyBytes(s)
	*pkb = *obj
	return err
}

// Creates an empty public key.
func NilPubkey() Pubkey {
	return Pubkey{}
}
