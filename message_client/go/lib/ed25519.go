package lib

import (
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"

	ccrypto "wraith.me/message_server/crypto"
)

const (
	ED25519_LEN = ccrypto.PUBKEY_SIZE
)

// Represents an Ed25519 keypair.
type Ed25519KP struct {
	SK          ccrypto.Privseed `json:"sk"`          //Holds the private key.
	PK          ccrypto.Pubkey   `json:"pk"`          //Holds the public key.
	Fingerprint string           `json:"fingerprint"` //Holds the fingerprint of the public key as a SHA-256 hash.
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

// Derives a `Privkey` object from this object.
func (kp Ed25519KP) Amalgamate() ccrypto.Privkey {
	bytes := [ccrypto.PRIVKEY_SIZE]byte{}
	copy(bytes[:ccrypto.PRIVKEY_SEED_SIZE], kp.SK[:])
	copy(bytes[ccrypto.PRIVKEY_SEED_SIZE:], kp.PK[:])
	return bytes
}

// Checks if this Ed25519 keypair is equal to another.
func (kp Ed25519KP) Equal(other Ed25519KP) bool {
	return subtle.ConstantTimeCompare(kp.SK[:], other.SK[:]) == 1 && subtle.ConstantTimeCompare(kp.PK[:], other.PK[:]) == 1
}

// Returns the JSON representation of the object.
func (kp Ed25519KP) JSON() string {
	json, _ := json.Marshal(kp)
	return string(json)
}

// Signs a message with this `Ed25519KP` object.
func (kp Ed25519KP) Sign(msg []byte) []byte {
	return ccrypto.Sign(kp.Amalgamate(), msg)
}

// Returns the string representation of the object.
func (kp Ed25519KP) String() string {
	return fmt.Sprintf("Ed25519KP{sk=%s, pk=%s}", hex.EncodeToString(kp.SK[:]), hex.EncodeToString(kp.PK[:]))
}

// Verifies a message and signature with this `Ed25519KP` object.
func (kp Ed25519KP) Verify(msg, sig []byte) bool {
	return ccrypto.Verify(kp.PK, msg, sig)
}
