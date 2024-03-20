//Adapted from: https://gist.github.com/saniales/532774ca61a17980431890fbef9438ad

package mongoutil

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// UUID represents a UUID as saved in MongoDB
type UUID struct{ uuid.UUID }

// NewUUID generates a new MongoDB compatible UUID.
func NewUUID4() (*UUID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &UUID{UUID: id}, nil
}

// Generates a new version 4 UUID. Panics if an error occurs.
func MustNewUUID4() *UUID {
	uuid, err := NewUUID4()
	if err != nil {
		panic(err)
	}
	return uuid
}

// NewUUID generates a new MongoDB compatible UUID.
func NewUUID7() (*UUID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	return &UUID{UUID: id}, nil
}

// Generates a new version 7 UUID. Panics if an error occurs.
func MustNewUUID7() *UUID {
	uuid, err := NewUUID7()
	if err != nil {
		panic(err)
	}
	return uuid
}

// UUIDFromStringOrNil returns a UUID parsed from the input string.
func UUIDFromStringOrNil(input string) *UUID {
	id := uuid.MustParse(input)
	if id == uuid.Nil {
		return nil
	}
	return &UUID{id}
}

// UUIDFromStringOrNil returns a UUID parsed from the input bytes.
func UUIDFromBytes(input []byte) *UUID {
	id := uuid.Must(uuid.FromBytes(input))
	return &UUID{id}
}

// MarshalBSONValue implements the bson.ValueMarshaler interface.
func (id UUID) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.TypeBinary, bsoncore.AppendBinary(nil, bson.TypeBinaryUUID, id.UUID[:]), nil
}

// UnmarshalBSONValue implements the bson.ValueUnmarshaler interface.
func (id *UUID) UnmarshalBSONValue(t bsontype.Type, raw []byte) error {
	//Ensure the incoming type is correct
	if t != bson.TypeBinary {
		return fmt.Errorf("(UUID) invalid format on unmarshalled bson value")
	}

	//Read the data from the BSON item
	_, data, _, ok := bsoncore.ReadBinary(raw)
	if !ok {
		return fmt.Errorf("(UUID) not enough bytes to unmarshal bson value")
	}
	copy(id.UUID[:], data)

	//No errors, so return nil
	return nil
}

// IsZero implements the bson.Zeroer interface.
func (id UUID) IsZero() bool {
	return bytes.Equal(id.UUID[:], uuid.Nil[:])
}

func (id UUID) Bytes() [16]byte {
	return id.UUID
}

func (id UUID) ShortString() string {
	return strings.ReplaceAll(id.String(), "-", "")
}
