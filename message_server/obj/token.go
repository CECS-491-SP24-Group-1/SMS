package obj

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"net"
	"time"

	"wraith.me/message_server/db/mongoutil"
)

//
//-- CLASS: UserToken
//

// Defines the length of the token's random bytes portion.
const RAND_TOKEN_LEN = 8

// Represents an opaque user login token.
type Token struct {
	//UserToken extends the abstract identifiable type.
	Identifiable `bson:",inline"`

	//The user that this token is for by ID.
	Subject mongoutil.UUID `json:"subject" bson:"subject"`

	//Defines the scope for which the token is allowed
	Scope TokenScope `json:"token_scope" bson:"token_scope"`

	//The IP address that created the token.
	CreationIP net.IP `json:"creation_ip" bson:"creation_ip"`

	//Denotes whether the token should expire
	Expire bool `json:"expire" bson:"expire"`

	//Denotes when the token will expire, as a Unix timestamp.
	Expiry int64 `json:"expiry" bson:"expiry"`

	//A array of random bytes. This field has no meaning on its own.
	//Rand [RAND_TOKEN_LEN]byte `json:"rand" bson:"token"`
}

func NewToken(subject mongoutil.UUID, creationIP net.IP, scope TokenScope, expiry time.Time) *Token {
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

func (ut *Token) FromB64(b64 string) error {
	//Decode to a byte array
	buf, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return err
	}

	//Decode from the byte array
	return ut.FromBytes(bytes.NewBuffer(buf))
}

func (ut *Token) FromBytes(buf *bytes.Buffer) error {
	denc := gob.NewDecoder(buf)
	return denc.Decode(ut)
}

func (ut Token) GetExpiry() time.Time {
	return time.Unix(ut.Expiry, 0)
}

func (ut Token) ToB64() string {
	buf, err := ut.ToBytes()
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

func (ut Token) ToBytes() (buf bytes.Buffer, err error) {
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(ut)
	return
}
