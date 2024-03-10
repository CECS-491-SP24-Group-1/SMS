package obj

import (
	"crypto/ed25519"
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

// Represents a user in the system.
type User struct {
	//The ID of the user.
	ID mongoutil.UUID `json:"id"`

	/*
		The username of the user. Can be changed at any time, but mustn't
		match that of another user. This field is case insensitive and must
		be 4-16 characters in length and only consist of alphanumeric characters
		and underscores.
	*/
	Username string `json:"username"`

	//The display name of the user. This must be 32 characters or less and is the username by default.
	DisplayName string `json:"display_name"`

	//The email of the user.
	Email string `json:"email"`

	//The user's public key. This must correspond to a private key held by the user.
	Pubkey [PUBKEY_SIZE]byte `json:"pubkey"`

	//The last time that the user logged in.
	LastLogin time.Time `json:"last_login"`

	//The last IP address that the user logged in from.
	LastIP net.IP `json:"last_ip"`

	//The user's global options, henceforth termed "user flags".
	Flags UserFlags `json:"flags"`
}

// Represents user options.
type UserFlags struct {
	//Whether the user's email has been verified.
	EmailVerified bool `json:"email_verified"`

	//Whether the user should be discoverable by their username.
	FindByUName bool `json:"find_by_uname"`

	//Controls who read receipts are sent to.
	ReadReceipts ReadReceiptsScope `json:"read_receipts"`

	//Whether the user can receive message from non-friended users.
	UnsolicitedMessages bool `json:"unsolicited_messages"`
}

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

// Unmarshals a read receipt flag to JSON.
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

// Controls the default flag options for new users.
func DefaultUserFlags() UserFlags {
	return UserFlags{
		EmailVerified:       false,   //Emails should be verified before user can do anything.
		FindByUName:         true,    //Users should be discoverable by their username by default.
		ReadReceipts:        FRIENDS, //Users should send read receipts only to their friends by default.
		UnsolicitedMessages: false,   //Users should not be able to be messaged without their consent by random, non-friends.
	}
}
