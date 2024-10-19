//go:generate go-enum --marshal --forceupper --mustparse --nocomments --names --values

package obj

//
//-- ENUM: IdType
//

/*
Defines what the ID corresponds to. An ID can correspond to users, servers,
messages, challenges, vaults, and so on.
*/
/*
ENUM(
	USER,
	SERVER,
	MESSAGE,
	CHALLENGE,
	TOKEN
)
*/
type IdType int8
