package obj

import "wraith.me/message_server/crypto"

//
//-- ABSTRACT CLASS: Entity
//

/*
Represents a generic entity in the system. This can either be a user or a
server. Each entity has an ID, type flag, and a public key.
*/
type Entity struct {
	//Entity extends the abstract identifiable type.
	Identifiable `bson:",inline"`

	//The entity's public key. This must correspond to a private key held by the entity.
	Pubkey crypto.PubkeyBytes `json:"pubkey" bson:"pubkey"`
}
