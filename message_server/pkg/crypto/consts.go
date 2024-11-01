package crypto

import "crypto/ed25519"

const (
	//The size of a private key in bytes.
	PRIVKEY_SIZE = ed25519.PrivateKeySize

	//The size of a private key's seed portion in bytes.
	PRIVKEY_SEED_SIZE = ed25519.SeedSize

	//The size of a private key's public portion in bytes.
	PRIVKEY_PUB_SIZE = ed25519.PublicKeySize

	//The size of a public key in bytes.
	PUBKEY_SIZE = PRIVKEY_PUB_SIZE

	//The size of a digital signature in bytes.
	SIG_SIZE = ed25519.SignatureSize

	//The size of a salt for Argon2id
	ARGON_SALT_SIZE = PRIVKEY_SEED_SIZE / 2
)
