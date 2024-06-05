package crypto

import (
	"crypto"
	"crypto/ed25519"
)

var (
	// The default algorithm used to hash signed messages.
	SignAlgo = crypto.SHA512
)

// Parses a byte slice to a `Privkey` or a `Pubkey`.
func MustFromBytes[K Privkey | Pubkey | Privseed | Signature](fun func([]byte) (K, error), bytes []byte) K {
	k, err := fun(bytes)
	if err != nil {
		panic(err)
	}
	return k
}

// Signs a message using a given `Privkey` object.
func Sign(privateKey Privkey, message []byte) Signature {
	sig := ed25519.Sign(privateKey[:], message)
	return MustFromBytes(SignatureFromBytes, sig)
}

// Verifies a message and signature using a given `Pubkey` object.
func Verify(publicKey Pubkey, message []byte, sig Signature) bool {
	return ed25519.Verify(publicKey[:], message, sig[:])
}
