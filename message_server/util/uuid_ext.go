package util

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

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
