//See: https://github.com/le5le-com/uuid

package mongoutil

import (
	"encoding/binary"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"wraith.me/message_server/pkg/util"
)

// Returns a MongoDB ObjectID from a given UUID. Preserves timestamp info to second precision.
// See: https://github.com/le5le-com/uuid/blob/0be35d1a92a89a0d2ffd56c5e2b0972fb1afeeee/uuid.go#L166
func UUID2OID(uuid util.UUID) primitive.ObjectID {
	//Allocate a byte array for the object ID; 12 bytes
	oid := primitive.ObjectID{}

	//Insert the Unix epoch into the OID; just the first 4 bytes in BE order
	//The time will be a garbage value if the UUID version doesn't support times (eg: UUIDv4)
	epochSec := uuid.Time().Second() //48-bit ts
	binary.BigEndian.PutUint32(oid[0:4], uint32(epochSec))

	//Copy the rest of the UUID into the OID
	uuidBytes := uuid.Bytes()
	oid[4] = uuidBytes[7]
	copy(oid[5:12], uuidBytes[9:16])

	//Return the OID
	return oid
}

// Returns a UUIDv7 from a given MongoDB ObjectID. Preserves timestamp info to second precision.
// See: https://github.com/le5le-com/uuid/blob/0be35d1a92a89a0d2ffd56c5e2b0972fb1afeeee/uuid.go#L113
func OID2UUID(oid primitive.ObjectID) util.UUID {
	//Allocate a byte array for the UUID; 16 bytes
	uuidBytes := [16]byte{}

	//Insert the Unix epoch into the OID; all bytes in BE order
	sec := int64(binary.BigEndian.Uint32(oid[0:4]))
	ms := time.Unix(sec, 0).UTC().UnixMilli()
	binary.BigEndian.PutUint64(uuidBytes[:], uint64(ms&((1<<48)-1)<<16))

	epochSec := oid.Timestamp().Unix() //32-bit ts
	binary.BigEndian.PutUint64(uuidBytes[0:8], uint64(epochSec))

	//Copy the rest of the OID into the UUID
	uuidBytes[6] = (uuidBytes[6] & 0x0f) | 0x70 // Version 7 [0111]
	uuidBytes[7] = oid[4]
	uuidBytes[8] = (uuidBytes[8] & 0x3f) | 0x80 // Variant [10]
	copy(uuidBytes[9:16], oid[5:])

	//Create a UUID from the byte array and return it
	return util.UUIDFromBytes(uuidBytes[:])
}
