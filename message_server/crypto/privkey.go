package crypto

import (
	"crypto"
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
)

const (
	//The size of a private key in bytes.
	PRIVKEY_SIZE = ed25519.PrivateKeySize

	//The size of a private key's seed portion in bytes.
	PRIVKEY_SEED_SIZE = ed25519.SeedSize

	//The size of a private key's public portion in bytes.
	PRIVKEY_PUB_SIZE = ed25519.PublicKeySize
)

//
//-- ALIAS: Privkey
//

// Represents the bytes of an entity's private key.
type Privkey [PRIVKEY_SIZE]byte

// Generates a new `Privkey` and `Pubkey` pair.
func NewKeypair(randSrc io.Reader) (Pubkey, Privkey, error) {
	pubkey, privkey, err := ed25519.GenerateKey(randSrc)
	if err != nil {
		return NilPubkey(), NilPrivkey(), err
	}
	return MustKeyFromBytes(PubkeyFromBytes, pubkey[:]), MustKeyFromBytes(PrivkeyFromBytes, privkey[:]), nil
}

// Creates an empty private key.
func NilPrivkey() Privkey {
	return Privkey{}
}

// Parses a `Privkey` object from a string.
func ParsePrivkeyBytes(str string) (Privkey, error) {
	//Derive a byte array from the string
	ba, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return NilPrivkey(), err
	}

	//Parse the resulting byte array
	return PrivkeyFromBytes(ba)
}

// Converts a byte slice into a new `Privkey` object.
func PrivkeyFromBytes(bytes []byte) (Privkey, error) {
	//Ensure proper length before parsing
	if len(bytes) != PRIVKEY_SEED_SIZE && len(bytes) != PRIVKEY_SIZE {
		return NilPrivkey(), fmt.Errorf("mismatched byte array size (%d); expected: %d or %d", len(bytes), PRIVKEY_SEED_SIZE, PRIVKEY_SIZE)
	}

	//If the size is only half that of a private key, do `scalar_mult` to derive the public bytes
	bin := [PRIVKEY_SIZE]byte{}
	if len(bytes) == PRIVKEY_SEED_SIZE {
		priv := ed25519.NewKeyFromSeed(bytes)
		copy(bin[:], priv)
	} else {
		copy(bin[:], bytes)
	}

	//Create a new object and return
	return Privkey(bin), nil
}

// Compares two `Privkey` objects.
func (prk Privkey) Equal(other Privkey) bool {
	return subtle.ConstantTimeCompare(prk[:], other[:]) == 1
}

// Gets the fingerprint of a `Privkey` object using SHA256.
func (prk Privkey) Fingerprint() string {
	hash := sha256.Sum256(prk[:])
	return hex.EncodeToString(hash[:])
}

// Marshals a `Privkey` object to JSON.
func (prk Privkey) MarshalJSON() ([]byte, error) {
	return json.Marshal(prk.String())
}

// Gets the public part of a `Privkey` object.
func (prk Privkey) Public() Pubkey {
	public := Pubkey{}
	copy(public[:], prk[PRIVKEY_SEED_SIZE:])
	return public
}

// Gets the private part of a `Privkey` object.
func (prk Privkey) Seed() Privseed {
	seed := Privseed{}
	copy(seed[:], prk[:PRIVKEY_SEED_SIZE])
	return seed
}

// Signs a message using this `Privkey` object.
func (prk Privkey) Sign(rand io.Reader, message []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	return ed25519.PrivateKey(prk[:]).Sign(rand, message, opts)
}

// Converts a `Privkey` object to a string.
func (prk Privkey) String() string {
	return base64.StdEncoding.EncodeToString(prk[:])
}

// Unmarshals a `Privkey` object from JSON.
func (prk *Privkey) UnmarshalJSON(b []byte) error {
	//Unmarshal to a string
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	//Derive a valid object from the string and reassign
	obj, err := ParsePrivkeyBytes(s)
	*prk = obj
	return err
}
