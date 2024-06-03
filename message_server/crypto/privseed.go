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
//-- ALIAS: Privseed
//

// Represents the "seed bytes" of an entity's private key.
type Privseed [PRIVKEY_SEED_SIZE]byte

// Creates an empty private key.
func NilPrivseed() Privseed {
	return Privseed{}
}

// Parses a `Privkey` object from a string.
func ParsePrivseedBytes(str string) (Privseed, error) {
	//Derive a byte array from the string
	ba, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return NilPrivseed(), err
	}

	//Parse the resulting byte array
	return PrivseedFromBytes(ba)
}

// Converts a byte slice into a new `Privseed` object.
func PrivseedFromBytes(bytes []byte) (Privseed, error) {
	//Ensure proper length before parsing
	if len(bytes) != PRIVKEY_SEED_SIZE {
		return NilPrivseed(), fmt.Errorf("mismatched byte array size (%d); expected: %d", len(bytes), PRIVKEY_SEED_SIZE)
	}

	//Create a new object and return
	return Privseed(bytes), nil
}

// Compares two `Privseed` objects.
func (prs Privseed) Equal(other Privseed) bool {
	return subtle.ConstantTimeCompare(prs[:], other[:]) == 1
}

// Gets the fingerprint of a `Privseed` object using SHA256.
func (prs Privseed) Fingerprint() string {
	hash := sha256.Sum256(prs[:])
	return hex.EncodeToString(hash[:])
}

// Marshals a `Privseed` object to JSON.
func (prs Privseed) MarshalJSON() ([]byte, error) {
	return json.Marshal(prs.String())
}

// Converts a `Privseed` object to a string.
func (prs Privseed) String() string {
	return base64.StdEncoding.EncodeToString(prs[:])
}

// Unmarshals a `Privseed` object from JSON.
func (prs *Privseed) UnmarshalJSON(b []byte) error {
	//Unmarshal to a string
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	//Derive a valid object from the string and reassign
	obj, err := ParsePrivseedBytes(s)
	*prs = obj
	return err
}
