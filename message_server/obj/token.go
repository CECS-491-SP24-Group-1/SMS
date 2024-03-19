package obj

import (
	"encoding/base64"
	"time"
	"unsafe" //The implications of using this package are understood.

	"wraith.me/message_server/db/mongoutil"
	"wraith.me/message_server/obj/ip_addr"
)

//
//-- CLASS: UserToken
//

const (
	// Defines the length of the token's random bytes portion.
	RAND_TOKEN_LEN = 8

	//Defines the size of the token in bytes
	TOKEN_SIZE_BYTES = int(unsafe.Sizeof(Token{}))
)

// Represents an opaque user login token.
type Token struct {
	//UserToken extends the abstract identifiable type.
	Identifiable `bson:",inline"`

	//The IP address that created the token.
	CreationIP ip_addr.IPAddr `json:"creation_ip" bson:"creation_ip"`

	//The user that this token is for by ID.
	Subject mongoutil.UUID `json:"subject" bson:"subject"`

	//Defines the scope for which the token is allowed
	Scope TokenScope `json:"token_scope" bson:"token_scope"`

	//Denotes whether the token should expire
	Expire bool `json:"expire" bson:"expire"`

	//Denotes when the token will expire, as a Unix timestamp.
	Expiry int64 `json:"expiry" bson:"expiry"`

	//A array of random bytes. This field has no meaning on its own.
	//Rand [RAND_TOKEN_LEN]byte `json:"rand" bson:"token"`
}

// Creates a new token object.
func NewToken(subject mongoutil.UUID, creationIP ip_addr.IPAddr, scope TokenScope, expiry time.Time) *Token {
	return &Token{
		Identifiable: Identifiable{
			ID:   *mongoutil.MustNewUUID7(),
			Type: IdTypeTOKEN,
		},
		Subject:    subject,
		Scope:      scope,
		CreationIP: creationIP,
		Expire:     true,
		Expiry:     expiry.Unix(),
		//Rand:       [8]byte(util.MustGenRandBytes(RAND_TOKEN_LEN)),
	}
}

// Creates a token from a base64 string.
func TokenFromB64(b64 string) (*Token, error) {
	buf, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}
	return TokenFromBytes(buf), nil
}

// Creates a token from a byte array. Thanks ChatGPT.
func TokenFromBytes(data []byte) *Token {
	var token Token
	sz := int(unsafe.Sizeof(token))
	if len(data) < sz {
		// Handle the case where the byte slice is smaller than the struct size
		// This could be due to an incomplete read or other issues
		// For simplicity, we'll return an empty struct
		return &token
	}
	// Interpret the byte slice as a pointer to the struct type
	return (*Token)(unsafe.Pointer(&data[0]))
}

// Checks if 2 tokens are equal.
func (ut Token) Equal(other Token) bool {
	return ut.ID == other.ID &&
		ut.Type == other.Type &&
		ut.Subject == other.Subject &&
		ut.Scope == other.Scope &&
		ut.CreationIP == other.CreationIP &&
		ut.Expire == other.Expire &&
		ut.Expiry == other.Expiry
}

// Gets the expiry time of the token.
func (ut Token) GetExpiry() time.Time {
	return time.Unix(ut.Expiry, 0)
}

// Converts a token into a string.
func (ut Token) String() string {
	return ut.ToB64()
}

// Converts the token into a base64 string.
func (ut Token) ToB64() string {
	return base64.StdEncoding.EncodeToString(ut.ToBytes())
}

// Converts a token into a byte array. See: https://stackoverflow.com/a/56272984
func (ut Token) ToBytes() []byte {
	return (*(*[TOKEN_SIZE_BYTES]byte)(unsafe.Pointer(&ut)))[:]
}
