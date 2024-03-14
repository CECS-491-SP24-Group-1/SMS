package obj

import (
	"encoding/json"
	"fmt"
	"strings"

	"wraith.me/message_server/db/mongoutil"
)

//
//-- CLASS: Challenge
//

/*
Represents any object in the system that is identifiable by a UUID. This
can include users, servers, messages, challenges, and so on.
*/
type Identifiable struct {
	//The ID of the object.
	ID mongoutil.UUID `json:"id" bson:"_id"`

	//The type of item this object is.
	Type IdType `json:"type" bson:"type"`
}

//
//-- ENUM: IdType
//

// Denotes whether an entity is a user or a server.
type IdType int

const (
	USER IdType = iota
	SERVER
)

// Converts an entity type flag to a string.
func (et IdType) String() string {
	ets := ""
	switch et {
	case USER:
		ets = "USER"
	case SERVER:
		ets = "SERVER"
	}
	return ets
}

// Converts an entity type flag string to an object.
func ParseIdType(s string) (IdType, error) {
	et := -1
	switch strings.ToUpper(s) {
	case "USER":
		et = int(USER)
	case "SERVER":
		et = int(SERVER)
	default:
		return -1, fmt.Errorf("IdType: invalid enum name '%s'", s)
	}
	return IdType(et), nil
}

// Marshals an entity type flag to JSON.
func (et IdType) MarshalJSON() ([]byte, error) {
	return json.Marshal(et.String())
}

// Unmarshals an entity type flag from JSON.
func (et *IdType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	if *et, err = ParseIdType(s); err != nil {
		return err
	}
	return nil
}
