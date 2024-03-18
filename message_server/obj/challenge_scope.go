//go:generate go-enum --marshal --forceupper --mustparse --nocomments --names --values
package obj

//
//-- ENUM: ChallengeScope
//

// Defines the environment or use case for this challenge.
/*
ENUM(
	EMAIL, //This type of challenge is issued to verify a user owns an email address.
	PUBKEY //This type of challenge is issued to verify a user owns a private key corresponding to the given public key.
)
*/
type ChallengeScope int
