package obj

import (
	"net"
	"time"

	"wraith.me/message_server/db/mongoutil"
	"wraith.me/message_server/util"
)

//TODO: add an equal function

//
//-- CLASS: User
//

// Represents a user of the system. A user is a type of entity.
type User struct {
	//User extends the abstract entity type.
	Entity `bson:",inline"`

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

	/*
		The user's tokens along with the IP address that created it. This should
		NOT be outputted if a JSON representation is requested.
	*/
	Tokens []UserToken `json:"-" bson:"tokens"`
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
			Identifiable: Identifiable{
				ID:   id,
				Type: IdTypeUSER,
			},
			Pubkey: pubkey,
		},
		Username:    username,
		DisplayName: displayName,
		Email:       email,
		LastLogin:   lastLogin,
		LastIP:      lastIP,
		Flags:       flags,
		Options:     options,
		Tokens:      make([]UserToken, 0),
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
		FindByUName:         true,                     //Users should be discoverable by their username by default.
		ReadReceipts:        ReadReceiptsScopeFRIENDS, //Users should send read receipts only to their friends by default.
		UnsolicitedMessages: false,                    //Users should not be able to be messaged without their consent by random, non-friends.
	}
}

//
//-- CLASS: UserToken
//

// Holds a pairing of a JWT token and various metrics about it.
type UserToken struct {
	//The IP address that created the token.
	CreationIP net.IP `json:"creation_ip" bson:"creation_ip"` //Might make sense to move to the token itself to prevent modification

	//The user agent of the device that created the token.
	CreationUA string `json:"creation_ua" bson:"creation_ua"` //Might make sense to move to the token itself to prevent modification

	//The last time that the token was used successfully by a user.
	LastUsage time.Time `json:"last_usage" bson:"last_usage"`

	//The contents of the token itself.
	Token string `json:"token" bson:"token"`
}
