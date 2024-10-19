package csolver

import (
	ccrypto "wraith.me/message_server/pkg/crypto"
	"wraith.me/message_server/pkg/util"
)

//TODO: move to the http_types/request

/*
Defines the structure of JSON form data sent in the 1st stage of a login
request. This contains the user's ID and public key, both of which must
match what's in the database.
*/
type LoginUser struct {
	//The UUID of the user to login as.
	ID util.UUID `json:"id" mapstructure:"id"`

	//The public key of the user to login as.
	PK ccrypto.Pubkey `json:"pk" mapstructure:"pk"`
}

/*
Defines the structure of JSON form data sent in the 2nd stage of a login
request. This contains everything that the 1st stage form data contains,
along with the token that was issued and the digital signature of the
token that was signed by the private key of the user.
*/
type LoginVerifyUser struct {
	//`loginVerifyUser` extends `loginUser` by adding the previously generated token and the client's signature.
	LoginUser `mapstructure:",squash"`

	//The login token that the user was given.
	Token string `json:"token" mapstructure:"token"`

	//The signature of the input token, signed by the user's private key.
	Signature ccrypto.Signature `json:"signature" mapstructure:"signature"`
}
