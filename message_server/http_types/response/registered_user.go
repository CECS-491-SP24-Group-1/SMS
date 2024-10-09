package response

import "wraith.me/message_server/util"

/*
Represents a user object that is passed to the client once registration
is completed.
*/
type RegisteredUser struct {
	//The ID of the user.
	ID util.UUID `json:"id"`

	//The username of the user.
	Username string `json:"username"`

	//The email of the user, but redacted.
	RedactedEmail string `json:"redacted_email"`

	//The fingerprint of the submitted public key.
	PKFingerprint string `json:"pk_fingerprint"`
}
