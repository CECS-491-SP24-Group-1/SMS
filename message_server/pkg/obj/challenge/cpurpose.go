//go:generate go-enum --marshal --forceupper --mustparse --nocomments --names --values
package challenge

//
//-- ENUM: CPurpose
//

//Defines the reason that this challenge was issued, eg: to create an
//account, to request account deletion, to confirm a user's identity,
//etc.
/*
ENUM(
	UNKNOWN, //The purpose of the challenge is unknown.
	REGISTER, //The purpose of the challenge is to complete account registration.
	LOGIN, //The purpose of the challenge is to perform account login.
	DELETE, //The purpose of the challenge is to complete account deletion.
	CONFIRM, //The purpose of the challenge is to confirm a claimed identity.
)
*/
type CPurpose int8
