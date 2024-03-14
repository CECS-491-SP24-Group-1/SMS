package obj

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"wraith.me/message_server/db/mongoutil"
	"wraith.me/message_server/util"
)

//TODO: add an equal function

//
//-- CLASS: User
//

// Represents a user in the system. A user is a type of entity.
type User struct {
	//User extends the abstract entity type.
	Entity

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

	//The last time that the user logged in.
	LastLogin time.Time `json:"last_login" bson:"last_login"`

	//The last IP address that the user logged in from.
	LastIP net.IP `json:"last_ip" bson:"last_ip"`

	//The user's flags. These mark items such as verification status, deletion, etc.
	Flags UserFlags `json:"flags" bson:"flags"`

	//The user's global options.
	Options UserOptions `json:"options" bson:"options"`
}

// Creates a new user object.
func NewUser(
	id mongoutil.UUID,
	pubkey PubkeyBytes,
	username string,
	displayName string,
	email string,
	lastLogin time.Time,
	lastIP net.IP,
	flags UserFlags,
	options UserOptions,
) *User {
	return &User{
		Entity: Entity{
			ID:     id,
			Type:   USER,
			Pubkey: pubkey,
		},
		Username:    username,
		DisplayName: displayName,
		Email:       email,
		LastLogin:   lastLogin,
		LastIP:      lastIP,
		Flags:       flags,
		Options:     options,
	}
}

// Creates a user from only a username and string. This should be used only for mocking.
func NewUserSimple(username string, email string) (*User, error) {
	//Precompute complex fields
	uuid, err := mongoutil.NewUUID7()
	if err != nil {
		return nil, err
	}
	randBytes, err := util.GenRandBytes(PUBKEY_SIZE)
	if err != nil {
		return nil, err
	}

	//Create the object
	return NewUser(
		*uuid,
		PubkeyBytes(randBytes),
		username,
		username,
		email,
		util.NowMillis(),
		net.ParseIP("127.0.0.1"),
		DefaultUserFlags(),
		DefaultUserOptions(),
	), nil
}

//
//-- CLASS: UserFlags
//

// Represents user flags.
type UserFlags struct {
	//Whether the user's email has been verified.
	EmailVerified bool `json:"email_verified" bson:"email_verified"`

	//Whether the user's public key has been verified to correspond to a private key.
	PubkeyVerified bool `json:"pubkey_verified" bson:"pubkey_verified"`

	//Whether the user's account has been marked for deletion. This flag is set to true by default, and is lifted when the 2 above flags are false.
	ShouldPurge bool `json:"should_purge" bson:"should_purge"`

	//The UTC time at which the account should be purged from the database. This field is ignored if `ShouldPurge` is false.
	PurgeBy time.Time `json:"purge_by" bson:"purge_by"`
}

// Controls the default flag options for new users.
func DefaultUserFlags() UserFlags {
	return UserFlags{
		EmailVerified:  false,                                //Emails should be verified before user can do anything.
		PubkeyVerified: false,                                //Public keys should be verified before user can do anything.
		ShouldPurge:    true,                                 //Accounts should be purged automatically by default due to missing verification of email and pubkey.
		PurgeBy:        util.NowMillis().Add(24 * time.Hour), //New accounts are purged after 24 hours by default if verification is not done.
	}
}

//
//-- CLASS: UserOptions
//

// Represents user options.
type UserOptions struct {
	//Whether the user should be discoverable by their username.
	FindByUName bool `json:"find_by_uname" bson:"find_by_uname"`

	//Controls who read receipts are sent to.
	ReadReceipts ReadReceiptsScope `json:"read_receipts" bson:"read_receipts"`

	//Whether the user can receive message from non-friended users.
	UnsolicitedMessages bool `json:"unsolicited_messages" bson:"unsolicited_messages"`
}

// Controls the default flag options for new users.
func DefaultUserOptions() UserOptions {
	return UserOptions{
		FindByUName:         true,    //Users should be discoverable by their username by default.
		ReadReceipts:        FRIENDS, //Users should send read receipts only to their friends by default.
		UnsolicitedMessages: false,   //Users should not be able to be messaged without their consent by random, non-friends.
	}
}

//
//-- ENUM: ReadReceiptsScope
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
