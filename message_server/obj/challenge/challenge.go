package challenge

import (
	"encoding/base64"
	"time"

	"wraith.me/message_server/db/mongoutil"
	"wraith.me/message_server/obj"
	"wraith.me/message_server/util"
)

const (
	//Defines how long a challenge should last. This is 30 minutes by default.
	DEFAULT_CHALLENGE_EXPIRY = time.Minute * 30

	//Defines the default size of the challenge payload.
	DEFAULT_CHALLENGE_PAYLOAD_SIZE = 24
)

var (
	// Defines the size of the challenge payload.
	ChallengePayloadSize = DEFAULT_CHALLENGE_PAYLOAD_SIZE

	//Defines the URL parameter name to target when getting a challenge solution.
	ChallengeURLParamName = "sol"
)

//
//-- CLASS: Challenge
//

/*
Represents a challenge given to a user to solve. A challenge can be used
to remove holds on accounts, prove identity, or provide authorization for
an account action such as deletion. A challenge can either be initiated by
a user or a server. Likewise, a challenge can either be responded to by a
user or a server, though the latter is not currently slated for immediate
implementation at this time.
*/
type Challenge struct {
	//Challenge extends the abstract identifiable type.
	obj.Identifiable `bson:",inline"`

	//The scope of the challenge
	Scope ChallengeScope `json:"scope" bson:"scope"`

	//The entity that initiated the challenge. This field may be abridged to save space.
	Initiator obj.Identifiable `json:"initiator" bson:"initiator"`

	//The entity that will respond to the challenge. This field may be abridged to save space.
	Responder obj.Identifiable `json:"responder" bson:"responder"`

	//The time at which the challenge will expire, irregardless of the status. This should be short to maximize security.
	Expiry time.Time `json:"expiry" bson:"expiry"`

	//The status of the challenge.
	Status ChallengeStatus `json:"status" bson:"status"`

	//The payload text of the challenge. This is what will actually be sent to a user.
	Payload string `json:"payload" bson:"payload"`
}

// Creates a new challenge, generating a random ID and challenge text for it.
func NewChallenge(
	scope ChallengeScope,
	initiator obj.Identifiable,
	responder obj.Identifiable,
	expiry time.Time,
) Challenge {
	return NewChallengeDeterministic(mongoutil.MustNewUUID7(), scope, initiator, responder, expiry)
}

// Creates a new challenge from an existing UUID.
func NewChallengeDeterministic(
	id mongoutil.UUID,
	scope ChallengeScope,
	initiator obj.Identifiable,
	responder obj.Identifiable,
	expiry time.Time,
) Challenge {
	//Create random challenge text
	ctext := base64.URLEncoding.EncodeToString(util.MustGenRandBytes(ChallengePayloadSize))

	//Create and return a challenge
	return Challenge{
		Identifiable: obj.Identifiable{
			ID:   id,
			Type: obj.IdTypeCHALLENGE,
		},
		Scope:     scope,
		Initiator: initiator,
		Responder: responder,
		Expiry:    expiry.Truncate(time.Millisecond).UTC(),
		Status:    ChallengeStatusPENDING,
		Payload:   ctext,
	}
}

// Checks if this challenge is equal to another.
func (ch Challenge) Equal(other Challenge) bool {
	return ch.ID == other.ID &&
		ch.Scope == other.Scope &&
		ch.Initiator.Equal(other.Initiator) &&
		ch.Responder.Equal(other.Responder) &&
		ch.Expiry == other.Expiry &&
		ch.Status == other.Status &&
		ch.Payload == other.Payload
}
