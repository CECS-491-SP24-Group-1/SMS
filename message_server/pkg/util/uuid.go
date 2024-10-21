//Adapted from: https://gist.github.com/saniales/532774ca61a17980431890fbef9438ad

package util

import (
	"bytes"
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// UUID represents a UUID as saved in MongoDB.
type UUID struct{ uuid.UUID }

// NewUUID generates a new MongoDB compatible UUID.
func NewUUID4() (res UUID, err error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return
	}
	return UUID{UUID: id}, nil
}

// Generates a new version 4 UUID. Panics if an error occurs.
func MustNewUUID4() UUID {
	uuid, err := NewUUID4()
	if err != nil {
		panic(err)
	}
	return uuid
}

// NewUUID generates a new MongoDB compatible UUID.
func NewUUID7() (res UUID, err error) {
	id, err := uuid.NewV7()
	if err != nil {
		return
	}
	return UUID{UUID: id}, nil
}

/*
NewUUID7FromTime generates a new MongoDB compatible UUID from an existing
time vs `time.Now()`.
See: https://www.perplexity.ai/search/generate-a-v7-uuid-from-a-give-K1LLrJpmR7Ode.4JFw6rPQ
*/
func NewUUID7FromTime(tim time.Time) UUID {
	//Get Unix timestamp in milliseconds
	ms := uint64(tim.UnixNano() / int64(time.Millisecond))

	//Create a 16-byte array for the UUID
	uuidBytes := [16]byte{}

	//Set the time_low, time_mid, and time_hi fields
	binary.BigEndian.PutUint32(uuidBytes[0:4], uint32(ms>>16))
	binary.BigEndian.PutUint16(uuidBytes[4:6], uint16(ms&0xFFFF))
	binary.BigEndian.PutUint16(uuidBytes[6:8], uint16(ms>>32))

	//Set the version (7)
	uuidBytes[6] = (uuidBytes[6] & 0x0F) | 0x70

	//Set the variant
	uuidBytes[8] = (uuidBytes[8] & 0x3F) | 0x80

	//Generate random bytes for the rest
	_, err := crand.Read(uuidBytes[9:])
	if err != nil {
		panic(fmt.Sprintf("v7New(); failed to generate random: %s", err))
	}

	//Return the full UUIDv7
	return UUID{uuidBytes}
}

// Generates a new version 7 UUID. Panics if an error occurs.
func MustNewUUID7() UUID {
	uuid, err := NewUUID7()
	if err != nil {
		panic(err)
	}
	return uuid
}

// Returns the "nil uuid": a UUID of all 0s.
func NilUUID() UUID {
	return UUID{[16]byte{}}
}

// Returns a UUID parsed from the input string.
func ParseUUIDv7(input string) (UUID, error) {
	//Parse the UUID
	id, err := uuid.Parse(input)
	if err != nil {
		return NilUUID(), err
	}

	//Ensure its a version 7
	if id.Version() != 7 {
		return NilUUID(), fmt.Errorf("must be a UUIDv7; got version %d", id.Version())
	}

	return UUID{id}, nil
}

/*
Returns a UUID parsed from the input string, or a nil UUID if the input
string is not a valid UUID.
*/
func UUIDFromString(input string) UUID {
	id := uuid.MustParse(input)
	if id == uuid.Nil {
		return NilUUID()
	}
	return UUID{id}
}

// Returns a UUID parsed from the input bytes.
func UUIDFromBytes(input []byte) UUID {
	id := uuid.Must(uuid.FromBytes(input))
	return UUID{id}
}

// Returns the bytes of the UUID.
func (id UUID) Bytes() [16]byte {
	return id.UUID
}

// Determines if a UUID is a nil UUID.
func (id UUID) IsNil() bool {
	return id.Bytes() == [16]byte{}
}

// IsZero implements the bson.Zeroer interface.
func (id UUID) IsZero() bool {
	return bytes.Equal(id.UUID[:], uuid.Nil[:])
}

// MarshalText implements the text marshaller method.
func (id UUID) MarshalText() ([]byte, error) {
	return []byte(id.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (id *UUID) UnmarshalText(text []byte) error {
	val := string(text)
	tmp, err := uuid.Parse(val)
	if err != nil {
		return err
	}
	*id = UUID{tmp}
	return nil
}

// Outputs a UUID with no separation hyphens.
func (id UUID) ShortString() string {
	return strings.ReplaceAll(id.String(), "-", "")
}

// Returns the time contained in the UUID.
func (id UUID) Time() time.Time {
	uuid := uuid.UUID(id.Bytes())
	return time.Unix(uuid.Time().UnixTime())
}

// Tests if a UUID is valid. This is shorthand for `uuid.Validate() == nil`.
func IsValidUUID(s string) bool {
	return uuid.Validate(s) == nil
}

// Tests if a UUIDv7 is valid.
func IsValidUUIDv7(s string) bool {
	//Test for validity before anything
	if !IsValidUUID(s) {
		return false
	}

	/*
		Parse the UUID into an object and check the version bit Using `MustParse()`
		is ok here since the UUID is guaranteed to be valid at this point
	*/
	return uuid.MustParse(s).Version() == 7
}
