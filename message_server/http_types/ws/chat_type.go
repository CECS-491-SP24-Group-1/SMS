//go:generate go-enum --marshal --forceupper --mustparse --nocomments --names --values

package obj

//
//-- ENUM: ChatType
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
type ChatType int8
