package crypto

import (
	"crypto/ed25519"
)

var (
// The default algorithm used to hash signed messages.
// SignAlgo = crypto.SHA512
)

// Parses a byte slice to a `Privkey` or a `Pubkey`.
func MustKeyFromBytes[K Privkey | Pubkey](fun func([]byte) (K, error), bytes []byte) K {
	k, err := fun(bytes)
	if err != nil {
		panic(err)
	}
	return k
}

// Signs a message using a given `Privkey` object.
func Sign(privateKey Privkey, message []byte) []byte {
	return ed25519.Sign(privateKey[:], message)
}

// Verifies a message using a given `Pubkey` object.
func Verify(publicKey Pubkey, message, sig []byte) bool {
	return ed25519.Verify(publicKey[:], message, sig)
}
