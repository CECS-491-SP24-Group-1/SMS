package obj

import (
	"encoding/base64"
	"fmt"
	"time"
	"unsafe" //The implications of using this package are understood.

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
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
	//Check if the data length is not the same as the expected size of a token
	if len(data) != TOKEN_SIZE_BYTES {
		//Return nil instead of a struct object
		return nil
	}
	//Interpret the byte slice as a pointer to the struct type and cast it to a token
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

// MarshalBSONValue implements the bson.ValueMarshaler interface.
func (ut Token) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(ut.ToB64())
	//return bson.TypeString, []byte(ut.ToB64()), nil
}

// Marshals a token to text. Used downstream by JSON and BSON marshalling.
func (ut Token) MarshalText() (text []byte, err error) {
	return []byte(ut.ToB64()), nil
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

// UnmarshalBSONValue implements the bson.ValueUnmarshaler interface.
func (ut *Token) UnmarshalBSONValue(t bsontype.Type, raw []byte) error {
	//Ensure the incoming type is correct
	if t != bson.TypeString {
		return fmt.Errorf("(Token) invalid format on unmarshalled bson value")
	}

	//Read the data from the BSON item
	var str string
	if err := bson.UnmarshalValue(bson.TypeString, raw, &str); err != nil {
		return err
	}

	//Deserialize the bytes into a struct
	obj, err := TokenFromB64(str)
	*ut = *obj
	return err
}

// Unmarshals a token from a string. Used downstream by JSON and BSON marshalling.
func (ut *Token) UnmarshalText(text []byte) error {
	tok, err := TokenFromB64(string(text[:]))
	*ut = *tok
	return err
}
