//go:generate go-enum --marshal --forceupper --mustparse --nocomments --names --values
package chat

//
//-- ENUM: Type
//

/*
Defines a Chat type.
*/
/*
ENUM(
	UKNOWN,
	EVENT,
	MESSAGE,
	EK,
	KEX1,
	KEX2,
)
*/
type Type int8
