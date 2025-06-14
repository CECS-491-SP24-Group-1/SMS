package db

import (
	"time"

	"github.com/qiniu/qmgo/field"
)

// Base struct for all storable DB objects.
type DBObj struct {
	//The ID of the object.
	//ID util.UUID `bson:"_id" json:"_id"` //TODO: eventually centralize this

	//The time at which the object was created.
	CreatedAt time.Time `bson:"created_at" json:"created_at"`

	//The time at which the object was last updated.
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// Creates a new DBObj. Call this when initializing the field in the struct.
func NewDBObj() DBObj {
	return DBObj{
		//ID: util.MustNewUUID7(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Adds the `CreatedAt` and `UpdatedAt` fields to objects persisted in the DB.
func (*DBObj) CustomFields() field.CustomFieldsBuilder {
	return field.NewCustom().SetCreateAt("CreatedAt").SetUpdateAt("UpdatedAt")
}
