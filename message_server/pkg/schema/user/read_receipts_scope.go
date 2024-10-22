//go:generate go-enum --marshal --forceupper --mustparse --nocomments --names --values

package user

//
//-- ENUM: ReadReceiptsScope
//

// Controls who read receipts are sent to.
/*
ENUM(
	EVERYONE	//Everyone is sent a read receipt.
	FRIENDS 	//Only friends are sent read receipts.
	NOBODY		//Nobody is sent a read receipt
)
*/
type ReadReceiptsScope int8
