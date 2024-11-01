package crypto

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

//
//-- ALIAS: Pubkey
//

// Represents the bytes of an entity's public key.
type Pubkey [PUBKEY_SIZE]byte

// Creates an empty public key.
func NilPubkey() Pubkey {
	return Pubkey{}
}

// Parses a `Pubkey` object from a string.
func ParsePubkey(str string) (Pubkey, error) {
	//Derive a byte array from the string
	ba, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return NilPubkey(), err
	}

	//Parse the resulting byte array
	return PubkeyFromBytes(ba)
}

// Converts a byte slice into a new `Pubkey` object.
func PubkeyFromBytes(bytes []byte) (Pubkey, error) {
	//Ensure proper length before parsing
	if len(bytes) != PUBKEY_SIZE {
		return NilPubkey(), fmt.Errorf("mismatched byte array size (%d); expected: %d", len(bytes), PUBKEY_SIZE)
	}

	//Create a new object and return
	bin := [PUBKEY_SIZE]byte{}
	subtle.ConstantTimeCopy(1, bin[:], bytes)

	//Create a new object and return
	return Pubkey(bin), nil
}

// Compares two `Pubkey` objects.
func (pkb Pubkey) Equal(other Pubkey) bool {
	return subtle.ConstantTimeCompare(pkb[:], other[:]) == 1
}

// Gets the fingerprint of a `Pubkey` object using SHA256.
func (pkb Pubkey) Fingerprint() string {
	hash := sha256.Sum256(pkb[:])
	return hex.EncodeToString(hash[:])
}

// Marshals a `Pubkey` object to JSON.
func (pkb Pubkey) MarshalJSON() ([]byte, error) {
	return json.Marshal(pkb.String())
}

// Marshals a `Pubkey` object to a string.
func (pkb Pubkey) MarshalText() ([]byte, error) {
	return []byte(pkb.String()), nil
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
	obj, err := ParsePubkey(s)
	*pkb = obj
	return err
}

// Unmarshals a `Pubkey` object from a string.
func (pkb *Pubkey) UnmarshalText(text []byte) error {
	var err error
	*pkb, err = ParsePubkey(string(text))
	return err
}

// Verifies a message and signature using this `Privkey` object.
func (pkb Pubkey) Verify(message []byte, sig Signature) bool {
	return Verify(pkb, message, sig)
}
