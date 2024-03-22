package obj

import (
	"wraith.me/message_server/db/mongoutil"
)

//
//-- CLASS: Identifiable
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
