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
//-- ALIAS: Signature
//

/*
Represents a digital signature, which is created by signing a message
with a `Privkey` object.
*/
type Signature [SIG_SIZE]byte

// Creates an empty signature.
func NilSignature() Signature {
	return Signature{}
}

// Parses a `Signature` object from a string.
func ParseSignature(str string) (Signature, error) {
	//Derive a byte array from the string
	ba, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return NilSignature(), err
	}

	//Parse the resulting byte array
	return SignatureFromBytes(ba)
}

// Converts a byte slice into a new `Signature` object.
func SignatureFromBytes(bytes []byte) (Signature, error) {
	//Ensure proper length before parsing
	if len(bytes) != SIG_SIZE {
		return NilSignature(), fmt.Errorf("mismatched byte array size (%d); expected: %d", len(bytes), SIG_SIZE)
	}

	//Create a new object and return
	return Signature(bytes), nil
}

// Compares two `Signature` objects.
func (sig Signature) Equal(other Signature) bool {
	return subtle.ConstantTimeCompare(sig[:], other[:]) == 1
}

// Gets the fingerprint of a `Signature` object using SHA256.
func (sig Signature) Fingerprint() string {
	hash := sha256.Sum256(sig[:])
	return hex.EncodeToString(hash[:])
}

// Marshals a `Signature` object to JSON.
func (sig Signature) MarshalJSON() ([]byte, error) {
	return json.Marshal(sig.String())
}

// Marshals a `Signature` object to a string.
func (sig Signature) MarshalText() ([]byte, error) {
	return []byte(sig.String()), nil
}

// Converts a `Signature` object to a string.
func (sig Signature) String() string {
	return base64.StdEncoding.EncodeToString(sig[:])
}

// Unmarshals a `Signature` object from JSON.
func (sig *Signature) UnmarshalJSON(b []byte) error {
	//Unmarshal to a string
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	//Derive a valid object from the string and reassign
	obj, err := ParseSignature(s)
	*sig = obj
	return err
}

// Unmarshals a `Signature` object from a string.
func (sig *Signature) UnmarshalText(text []byte) error {
	var err error
	*sig, err = ParseSignature(string(text))
	return err
}
