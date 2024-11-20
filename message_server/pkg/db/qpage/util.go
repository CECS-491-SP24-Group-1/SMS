package qpage

import (
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Asserts that the datatype of the destination array is ok.
func assertCorrectOutputType(dest any) (reflect.Value, error) {
	//Get the type of the destination slice
	destValue := reflect.ValueOf(dest)

	//Ensure dest is a pointer to a slice
	if destValue.Kind() != reflect.Ptr || destValue.Elem().Kind() != reflect.Slice {
		return reflect.Value{}, fmt.Errorf("dest must be a pointer to a slice")
	}

	return destValue, nil
}

// Gets the value from a BSON D given a key.
func getValueFromBsonD(d bson.D, key string) (interface{}, bool) {
	for _, elem := range d {
		if elem.Key == key {
			return elem.Value, true
		}
	}
	return nil, false
}

// Gets the value of an ID from a BSON document.
func idString(id interface{}) string {
	switch val := id.(type) {
	case primitive.ObjectID:
		return val.Hex()
	case string:
		return val

	//This section will be hit for `UUID` types
	case primitive.Binary:
		//UUID (old) or UUID (https://bsonspec.org/spec.html)
		if val.Subtype == 3 || val.Subtype == 4 {
			uuid, err := uuid.FromBytes(val.Data)
			if err == nil {
				return uuid.String()
			}
		}
		return fmt.Sprintf("%v", val)

	//Default fall-through: simply `sprintf` the value
	default:
		return fmt.Sprintf("%v", val)
	}
}

// Unmarshals a BSON document to a reflected type.
func unmarshalBsonD(data bson.D, targetType reflect.Type) (interface{}, error) {
	//Create a new value of the target type
	value := reflect.New(targetType).Interface()

	//Convert bson.D to raw BSON bytes
	bsonBytes, err := bson.Marshal(data)
	if err != nil {
		return nil, err
	}

	//Unmarshal the BSON bytes into the new value
	err = bson.Unmarshal(bsonBytes, value)
	if err != nil {
		return nil, err
	}

	//Return the dereferenced value
	return reflect.ValueOf(value).Elem().Interface(), nil
}
