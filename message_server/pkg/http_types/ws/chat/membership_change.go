package chat

import "encoding/json"

// Represents a membership change event.
type MembershipChange struct {
	//The number of users in the room before the update.
	Old int `json:"old"`

	//The number of users in the room after the update.
	New int `json:"new"`
}

// Marshals the message to JSON.
func (m MembershipChange) JSON() []byte {
	jsons, err := json.Marshal(m)
	if err != nil {
		panic("Message::JSON: " + err.Error())
	}
	return jsons
}
