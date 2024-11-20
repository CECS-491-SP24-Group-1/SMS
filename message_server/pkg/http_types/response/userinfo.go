package response

import (
	"wraith.me/message_server/pkg/crypto"
	"wraith.me/message_server/pkg/schema/user"
	"wraith.me/message_server/pkg/util"
)

// Represents a user info object returned from looking up a user.
type UInfo struct {
	ID          util.UUID     `json:"id" bson:"_id"`
	Pubkey      crypto.Pubkey `json:"pubkey" bson:"pubkey"`
	Username    string        `json:"username" bson:"username"`
	DisplayName string        `json:"display_name" bson:"display_name"`
}

// Emits a user info object from an existing user.
func FromUser(user user.User) UInfo {
	return UInfo{
		ID:          user.ID,
		Pubkey:      user.Pubkey,
		Username:    user.Username,
		DisplayName: user.DisplayName,
	}
}
