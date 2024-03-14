package obj

import (
	"encoding/json"
	"fmt"
	"strings"

	"wraith.me/message_server/db/mongoutil"
)

//
//-- ABSTRACT CLASS: Entity
//

/*
Represents an entity in the system. This can either be a user or a server.
Each entity has an ID, designation, and a public key.
*/
type Entity struct {
	//The ID of the entity.
	ID mongoutil.UUID `json:"id" bson:"_id"`

	//The type of entity this object is.
	Type EntityType `json:"type" bson:"type"`

	//The entity's public key. This must correspond to a private key held by the entity.
	Pubkey PubkeyBytes `json:"pubkey" bson:"pubkey"`
}

//
//-- ENUM: EntityType
//

// Denotes whether an entity is a user or a server.
type EntityType int

const (
	USER EntityType = iota
	SERVER
)

// Converts an entity type flag to a string.
func (et EntityType) String() string {
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
func ParseEntityType(s string) (EntityType, error) {
	et := -1
	switch strings.ToUpper(s) {
	case "USER":
		et = int(USER)
	case "SERVER":
		et = int(SERVER)
	default:
		return -1, fmt.Errorf("EntityType: invalid enum name '%s'", s)
	}
	return EntityType(et), nil
}

// Marshals an entity type flag to JSON.
func (et EntityType) MarshalJSON() ([]byte, error) {
	return json.Marshal(et.String())
}

// Unmarshals an entity type flag from JSON.
func (et *EntityType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	if *et, err = ParseEntityType(s); err != nil {
		return err
	}
	return nil
}
