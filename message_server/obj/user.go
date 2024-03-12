package obj

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"wraith.me/message_server/db/mongoutil"
)

const (
	PUBKEY_SIZE = ed25519.PublicKeySize
)

//
//-- CLASS: User
//

// Represents a user in the system.
type User struct {
	//The ID of the user.
	ID mongoutil.UUID `json:"id" bson:"_id"`

	/*
		The username of the user. Can be changed at any time, but mustn't
		match that of another user. This field is case insensitive, and must
		be 4-16 characters in length and only consist of alphanumeric characters
		and underscores.
	*/
	Username string `json:"username" bson:"username"`

	//The display name of the user. This must be 32 characters or less and is the username by default.
	DisplayName string `json:"display_name" bson:"display_name"`

	//The email of the user.
	Email string `json:"email" bson:"email"`

	//The user's public key. This must correspond to a private key held by the user.
	Pubkey PubkeyBytes `json:"pubkey" bson:"pubkey"`

	//The last time that the user logged in.
	LastLogin time.Time `json:"last_login" bson:"last_login"`

	//The last IP address that the user logged in from.
	LastIP net.IP `json:"last_ip" bson:"last_ip"`

	//The user's global options, henceforth termed "user flags".
	Flags UserFlags `json:"flags" bson:"flags"`
}

//
//-- CLASS: PubkeyBytes
//

// Represents the bytes of the user's public key.
type PubkeyBytes [PUBKEY_SIZE]byte

// Marshals a `PubkeyBytes` object to JSON.
func (pkb PubkeyBytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(pkb.String())
}

// Parses a `PubkeyBytes` object from a string.
func ParsePubkeyBytes(str string) (*PubkeyBytes, error) {
	//Derive a byte array from the string
	ba, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}

	//Ensure the byte array length is correct
	if len(ba) != PUBKEY_SIZE {
		return nil, fmt.Errorf("mismatched byte array size (%d); expected: %d", len(ba), PUBKEY_SIZE)
	}

	//Copy the bytes to a new object and return it
	obj := &PubkeyBytes{}
	copy(obj[:], ba)
	return obj, nil
}

// Converts a `PubkeyBytes` object to a string.
func (pkb PubkeyBytes) String() string {
	return base64.StdEncoding.EncodeToString(pkb[:])
}

// Unmarshals a `PubkeyBytes` object from JSON.
func (pkb *PubkeyBytes) UnmarshalJSON(b []byte) error {
	//Unmarshal to a string
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	//Derive a valid object from the string and reassign
	obj, err := ParsePubkeyBytes(s)
	*pkb = *obj
	return err
}

//
//-- CLASS: UserFlags
//

// Represents user options.
type UserFlags struct {
	//Whether the user's email has been verified.
	EmailVerified bool `json:"email_verified" bson:"email_verified"`

	//Whether the user should be discoverable by their username.
	FindByUName bool `json:"find_by_uname" bson:"find_by_uname"`

	//Controls who read receipts are sent to.
	ReadReceipts ReadReceiptsScope `json:"read_receipts" bson:"read_receipts"`

	//Whether the user can receive message from non-friended users.
	UnsolicitedMessages bool `json:"unsolicited_messages" bson:"unsolicited_messages"`
}

// Controls the default flag options for new users.
func DefaultUserFlags() UserFlags {
	return UserFlags{
		EmailVerified:       false,   //Emails should be verified before user can do anything.
		FindByUName:         true,    //Users should be discoverable by their username by default.
		ReadReceipts:        FRIENDS, //Users should send read receipts only to their friends by default.
		UnsolicitedMessages: false,   //Users should not be able to be messaged without their consent by random, non-friends.
	}
}

//
//-- CLASS: ReadReceiptsScope
//

// Controls who read receipts are sent to.
type ReadReceiptsScope int

const (
	EVERYONE ReadReceiptsScope = iota //Everyone is sent a read receipt.
	FRIENDS                           //Only friends are sent read receipts.
	NOBODY                            //Nobody is sent a read receipt
)

// Converts a read receipt flag to a string.
func (rr ReadReceiptsScope) String() string {
	rrs := ""
	switch rr {
	case EVERYONE:
		rrs = "EVERYONE"
	case FRIENDS:
		rrs = "FRIENDS"
	case NOBODY:
		rrs = "NOBODY"
	}
	return rrs
}

// Converts a read receipt flag string to an object.
func ParseReadReceiptsScope(s string) (ReadReceiptsScope, error) {
	rri := -1
	switch strings.ToUpper(s) {
	case "EVERYONE":
		rri = int(EVERYONE)
	case "FRIENDS":
		rri = int(FRIENDS)
	case "NOBODY":
		rri = int(NOBODY)
	default:
		return -1, fmt.Errorf("ReadReceiptsScope: invalid enum name '%s'", s)
	}
	return ReadReceiptsScope(rri), nil
}

// Marshals a read receipt flag to JSON.
func (rr ReadReceiptsScope) MarshalJSON() ([]byte, error) {
	return json.Marshal(rr.String())
}

// Unmarshals a read receipt flag from JSON.
func (rr *ReadReceiptsScope) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	if *rr, err = ParseReadReceiptsScope(s); err != nil {
		return err
	}
	return nil
}
