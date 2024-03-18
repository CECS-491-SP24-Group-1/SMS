package challenge

import (
	"time"

	"wraith.me/message_server/db/mongoutil"
	"wraith.me/message_server/obj"
)

const (
	//Defines how long a challenge should last. This is 30 minutes by default.
	DEFAULT_CHALLENGE_EXPIRY = time.Minute * 30
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
}

func NewChallenge(
	scope ChallengeScope,
	initiator obj.Identifiable,
	responder obj.Identifiable,
) *Challenge {
	return &Challenge{
		Identifiable: obj.Identifiable{
			ID:   *mongoutil.MustNewUUID7(),
			Type: obj.IdTypeCHALLENGE,
		},
		Scope:     scope,
		Initiator: initiator,
		Responder: responder,
		Expiry:    time.Now().UTC().Add(DEFAULT_CHALLENGE_EXPIRY),
		Status:    ChallengeStatusPENDING,
	}
}

//TODO: extend challenge by adding challenge text and solution
