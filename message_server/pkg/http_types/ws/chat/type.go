//go:generate go-enum --marshal --forceupper --mustparse --nocomments --names --values
package chat

//
//-- ENUM: Type
//

// Defines a Chat type.
/*
ENUM(
	UNKNOWN		//Unknown chat type.
	U_MSG		//User message.
	S_MSG		//Server message.
	S_ERR		//Server error message.
	JOIN_EVENT	//A user joined the room.
	QUIT_EVENT	//A user left the room.
	MEMBERSHIP	//Membership announcement message.
	EK			//An encryption key sent by a user for the purpose of decrypting a group message.
	KEX1		//Step 1 of an X3DH KEX operation.
	KEX2		//Step 2 of an X3DH KEX operation.
)
*/
type Type int8
