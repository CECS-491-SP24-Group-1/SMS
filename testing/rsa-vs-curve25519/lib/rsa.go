package lib

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// Represents an RSA keypair.
type RSAKP struct {
	SK          []byte `json:"sk"`          //Holds the private key.
	PK          []byte `json:"pk"`          //Holds the public key.
	Fingerprint string `json:"fingerprint"` //Holds the fingerprint of the public key as a SHA-256 hash.
}

// Generates a new RSA keypair.
func RSAKeygen(size int) RSAKP {
	//Ensure the size is valid for RSA
	if size != 1024 && size != 2048 && size != 3072 && size != 4096 {
		panic("RSA key size must be 1024, 2048, 3072, or 4096")
	}

	//Generate the RSA key
	key, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		panic(err)
	}

	//Marshal into a byte array and create the object
	return RSAFromBytes(x509.MarshalPKCS1PrivateKey(key), x509.MarshalPKCS1PublicKey(&key.PublicKey))
}

// Derives an RSA keypair object from raw bytes.
func RSAFromBytes(sk []byte, pk []byte) RSAKP {
	//Create the base object
	out := RSAKP{}

	//Assign the public and private key bytes
	out.SK = make([]byte, len(sk))
	out.PK = make([]byte, len(pk))
	copy(out.SK[:], sk[:])
	copy(out.PK[:], pk[:])

	//Hash the public key
	h := sha256.Sum256(pk)
	out.Fingerprint = hex.EncodeToString(h[:])

	//Return the object
	return out
}

// Returns the JSON representation of the object.
func (kp RSAKP) JSON() string {
	json, _ := json.Marshal(kp)
	return string(json)
}

// Returns the string representation of the object.
func (kp RSAKP) String() string {
	return fmt.Sprintf("RSAKP{sk=%s, pk=%s}", hex.EncodeToString(kp.SK[:]), hex.EncodeToString(kp.PK[:]))
}
