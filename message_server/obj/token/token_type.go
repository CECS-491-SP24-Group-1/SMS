//go:generate go-enum --marshal --forceupper --mustparse --nocomments --names --values
package token

//
//-- ENUM: TokenType
//

/*
Defines the type of token this is, whether it be auth or refresh.
*/
/*
ENUM(
	NONE = 		0
	ACCESS = 	1
	REFRESH = 	2
)
*/
type TokenType uint8
