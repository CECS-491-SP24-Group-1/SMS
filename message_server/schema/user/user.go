package user

import (
	"crypto/subtle"
	"net"
	"time"

	"wraith.me/message_server/crypto"
	"wraith.me/message_server/obj"
	"wraith.me/message_server/obj/ip_addr"
	"wraith.me/message_server/util"
)

//TODO: add an equal function

//
//-- CLASS: User
//

// Represents a user of the system. A user is a type of entity.
type User struct {
	//User extends the abstract entity type.
	obj.Entity `json:",inline" bson:",inline,squash"`

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
	LastIP ip_addr.IPAddr `json:"last_ip" bson:"last_ip"`

	//The user's flags. These mark items such as verification status, deletion, etc.
	Flags UserFlags `json:"flags" bson:"flags"`

	//The user's global options.
	Options UserOptions `json:"options" bson:"options"`

	//The user's refresh tokens, keyed by their IDs in string form.
	Tokens map[string]UserToken `json:"tokens" bson:"tokens"`
}

//-- Constructors

// Creates a new user object.
func NewUser(
	id util.UUID,
	pubkey crypto.Pubkey,
	username string,
	displayName string,
	email string,
	lastLogin time.Time,
	lastIP ip_addr.IPAddr,
	flags UserFlags,
	options UserOptions,
) *User {
	return &User{
		Entity: obj.Entity{
			Identifiable: obj.Identifiable{
				ID:   id,
				Type: obj.IdTypeUSER,
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
		Tokens:      make(map[string]UserToken, 0),
	}
}

// Creates a user from only a username and string. This should be used only for mocking.
func NewUserSimple(username string, email string) *User {
	return NewUser(
		util.MustNewUUID7(),
		crypto.Pubkey(util.Must(util.GenRandBytes(crypto.PUBKEY_SIZE))),
		username,
		username,
		email,
		util.NowMillis(),
		ip_addr.FromNetIP(net.ParseIP("127.0.0.1")),
		DefaultUserFlags(),
		DefaultUserOptions(),
	)
}

//-- Methods

// Adds a new refresh token to this user object.
func (u *User) AddToken(tid, token string, exp time.Time) {
	//Create the token map if it doesn't already exist
	if u.Tokens == nil {
		u.Tokens = make(map[string]UserToken)
	}

	//Add the token to the list of the user's tokens
	tok := UserToken{Token: token, Expiry: exp}
	(*u).Tokens[tid] = tok
}

// Checks if a user has a particular token.
func (u User) HasToken(tok string) bool {
	for _, token := range u.Tokens {
		if subtle.ConstantTimeCompare([]byte(token.Token), []byte(tok)) == 1 {
			return true
		}
	}
	return false
}

// Checks if a user has a particular token by ID.
func (u User) HasTokenById(tokId string) bool {
	_, ok := u.Tokens[tokId]
	return ok
}

// Marks a user's email as verified.
func (u *User) MarkEmailVerified() {
	u.Flags.EmailVerified = true
	if u.Flags.EmailVerified && u.Flags.PubkeyVerified {
		u.Flags.ShouldPurge = false
	}
}

// Marks a user's public key as verified.
func (u *User) MarkPKVerified() {
	u.Flags.PubkeyVerified = true
	if u.Flags.EmailVerified && u.Flags.PubkeyVerified {
		u.Flags.ShouldPurge = false
	}
}

// Removes a refresh token from this user object.
func (u *User) RemoveToken(tid string) {
	delete(u.Tokens, tid)
}

// Unmarks a user's email as verified.
func (u *User) UnmarkEmailVerified() {
	u.Flags.EmailVerified = false
	if !u.Flags.ShouldPurge {
		u.Flags.ShouldPurge = true
	}
}

// Unmarks a user's public key as verified.
func (u *User) UnmarkPKVerified() {
	u.Flags.PubkeyVerified = false
	if !u.Flags.ShouldPurge {
		u.Flags.ShouldPurge = true
	}
}

//-- Embedded class definitions

//
//-- CLASS: UserFlags
//

// Represents user flags.
type UserFlags struct {
	//Indicates if the user's email has been verified.
	EmailVerified bool `json:"email_verified" bson:"email_verified"`

	//Indicates if the user's public key has been verified to correspond to a private key.
	PubkeyVerified bool `json:"pubkey_verified" bson:"pubkey_verified"`

	//Indicates if the user's account has been marked for deletion. This flag is set to true by default, and is lifted when the 2 above flags are false.
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
	//Indicates if the user should be discoverable by their username.
	FindByUName bool `json:"find_by_uname" bson:"find_by_uname"`

	//Controls who read receipts are sent to.
	ReadReceipts ReadReceiptsScope `json:"read_receipts" bson:"read_receipts"`

	//Indicates if the user can receive messages from non-friended users.
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

// Represents a single refresh token.
type UserToken struct {
	//The token itself.
	Token string `json:"token" bson:"token"`

	//The expiry of the token.
	Expiry time.Time `json:"expiry" bson:"expiry"`
}
