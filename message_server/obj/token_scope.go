//go:generate go-enum --marshal --forceupper --mustparse --nocomments --names --values
package obj

//
//-- ENUM: TokenScope
//

// Defines the scope for which the token is valid
/*
ENUM(
	POST_SIGNUP = 1 //The token is only allowed to complete the login challenges.
	USER //The token is allowed to be used everywhere that a normal user can access.
)
*/
type TokenScope uint8
