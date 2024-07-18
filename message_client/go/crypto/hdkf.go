package crypto

import (
	"crypto/sha256"
	"fmt"
	"hash"
	"io"

	"golang.org/x/crypto/hkdf"
)

// Sets the minimum the salt and destination buffers must be for an HKDF operation.
const _HKDF_MIN_BUF_SIZE = 8

var (
	//Sets the hash function for the HKDF operation. Uses SHA-2 256 by default.
	HKDFHashFunc func() hash.Hash = sha256.New
)

// Derives an Ed25519 key from a passphrase, salt, and optional info.
func Ed25519HKDF(passphrase string, salt, info []byte) (Ed25519KP, error) {
	privkey := make([]byte, 32)
	err := hkdfHelper(&privkey, passphrase, salt, info)
	return Ed25519FromSK(privkey), err
}

// Derives x random bytes from a passphrase, salt, and optional info.
func hkdfHelper(dest *[]byte, passphrase string, salt, info []byte) error {
	//Ensure the destination pointer can hold data
	if len(*dest) < _HKDF_MIN_BUF_SIZE || len(salt) < _HKDF_MIN_BUF_SIZE {
		return fmt.Errorf("destination or salt buffers too small; minimum size is %d bytes", _HKDF_MIN_BUF_SIZE)
	}

	//Create a new HKDF instance and read x bytes into the destination buffer
	hkdfReader := hkdf.New(HKDFHashFunc, []byte(passphrase), salt, info)
	if _, err := io.ReadFull(hkdfReader, *dest); err != nil {
		return err
	}
	return nil
}
