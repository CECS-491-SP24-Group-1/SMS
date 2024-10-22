//go:generate go-enum --marshal --forceupper --mustparse --nocomments --names --values
package challenge

//
//-- ENUM: CType
//

// Defines the type of challenge this is.
/*
ENUM(
	UNKNOWN 	//The type of the challenge is unknown.
	EMAIL 		//This type of challenge is issued to verify that a user owns an email address.
	PUBKEY		//This type of challenge is issued to verify that a user owns a private key corresponding to the given public key.
)
*/
type CType int8
