//go:generate go-enum --marshal --forceupper --mustparse --nocomments --names --values
package notification

//
//-- ENUM: Type
//

// Defines the type of notification this is.
/*
ENUM(
	UNKNOWN		//The type of the notification is unknown.
	NEW_MSG		//A notification fired off when a new message is received.
	FRQ_ACCEPT	//A notification fired off when a friend request was accepted.
	FRQ_REJECT	//A notification fired off when a friend request was rejected.
	FRQ_NEW		//A notification fired off when a friend request has been received.
)
*/
type Type int8
