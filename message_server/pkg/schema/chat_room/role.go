//go:generate go-enum --marshal --forceupper --mustparse --nocomments --names --values

package chatroom

//
//-- ENUM: Role
//

// Sets the roles that users can have in a chatroom.
/*
ENUM(
	MEMBER 		//The user is a regular member of the room.
	MODERATOR 	//The user can add and remove members
	OWNER 		//The user is the owner of the group and enjoys the same rights as moderators, but can delete the group too.
)
*/
type Role int8
