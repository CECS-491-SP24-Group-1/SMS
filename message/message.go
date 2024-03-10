package message

import (
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
)

//--- GENERIC MESSAGE START

/* Describes a generic message that can be exchanged between 2+ parties. */
type GenericMessage struct {
	//The ID of the message, represented by a version 7 UUID.
	ID uuid.UUID `json:"id"`
}

// --Constructors
/** Creates a new generic message. */
func NewGenericMessage() *GenericMessage {
	id, _ := uuid.NewV7()
	return &GenericMessage{
		ID: id,
	}
}

//--Methods
/* Returns the time the message was created in terms of the Unix epoch. This is encoded inside the UUID. */
func (m GenericMessage) Created() time.Time {
	//Ensure the UUID is a version 7 UUID
	if m.ID.Version() != 7 {
		return time.UnixMicro(0) //Return Jan 1, 1960
	}

	//Get the Unix time from the UUID
	usec, unsec := m.ID.Time().UnixTime()
	return time.Unix(usec, unsec)
}

/* Returns the string representation of the message. */
func (m GenericMessage) String() string {
	return fmt.Sprintf(
		"GenericMessage{id=%s}",
		m.ID.String(),
	)
}

//--- GENERIC MESSAGE END

//--- EXPIRING MESSAGE START

/* Describes a message that expires after a given time. */
type ExpiringMessage struct {
	GenericMessage

	Expiry time.Time `json:"expiry"` //The time at which the message will expire.
}

// --Constructors
/** Creates a new expiring message. This function accepts either a Time or Duration. */
func NewExpiringMessage(expiry interface{}) (*ExpiringMessage, error) {
	//Switch over the expiry and get the correct type
	var parsedExpiry time.Time
	switch v := expiry.(type) {
	case time.Time:
		parsedExpiry = v
	case time.Duration:
		parsedExpiry = time.Now().Add(v)
	default:
		return nil, fmt.Errorf("unrecognized duration type: %s", reflect.TypeOf(expiry))
	}

	return &ExpiringMessage{
		GenericMessage: *NewGenericMessage(),
		Expiry:         parsedExpiry,
	}, nil
}

//--Methods
/* Returns how long the message has left before it expires. Returns a zero duration if it's already expired. */
func (m ExpiringMessage) DurationToExpiry() time.Duration {
	//Check if the message has expired, otherwise return the duration left
	if m.IsExpired() {
		return time.Duration(0)
	}
	return time.Until(m.Expiry)
}

/* Marks the message as expired. */
func (m *ExpiringMessage) ExpireNow() {
	m.Expiry = time.Now()
}

/* Determines whether the message will be expired when a given time comes to pass. */
func (m ExpiringMessage) IsExpiredAt(time time.Time) bool {
	return time.After(m.Expiry)
}

/* Determines whether the message is expired at the current time. */
func (m ExpiringMessage) IsExpired() bool {
	return m.IsExpiredAt(time.Now())
}

/* Returns the string representation of the message. */
func (m ExpiringMessage) String() string {
	return fmt.Sprintf(
		"ExpiringMessage{id=%s, expiry=%d}",
		m.ID.String(),
		m.Expiry.Unix(),
	)
}

//--- EXPIRING MESSAGE END
