package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"wraith.me/message_server/db/mongoutil"
)

func TestUUIDvsOID(t *testing.T) {
	ts := time.Now()
	uuid, _ := uuid.NewV7()
	oid := primitive.NewObjectID()

	fmt.Printf("ts:\t\t%s\n", ts.Format(time.RFC3339Nano))
	fmt.Printf("uuid:\t\t%s\n", uuid)
	fmt.Printf("oid:\t\t%s\n", oid.Hex())
	fmt.Printf("ts (uuid):\t%s\n", time.Unix(uuid.Time().UnixTime()).Format(time.RFC3339Nano))
	fmt.Printf("ts (oid):\t%s\n", oid.Timestamp().Format(time.RFC3339Nano))
}

func TestUUIDxOID(t *testing.T) {
	//Generate UUID and get the timestamp
	uuidIn, _ := mongoutil.NewUUID7()
	uuidiTS := time.Unix(uuidIn.Time().UnixTime()).UTC()

	//Print stuff
	fmt.Printf("UUID In:  %s\n", uuidIn.String())
	fmt.Printf("TS In:\t  %s -> %d\n", uuidiTS.Format(time.RFC3339Nano), uuidiTS.Unix())

	//Convert the UUID to an OID and get the timestamp
	oid := mongoutil.UUID2OID(uuidIn)
	oidTS := oid.Timestamp()

	//Print stuff
	fmt.Printf("OID:\t  %s\n", oid.Hex())
	fmt.Printf("TS:\t  %s -> %d\n", oidTS.Format(time.RFC3339Nano), oidTS.Unix())

	//Convert the OID to a UUIDv7 and get the timestamp
	uuidOut := mongoutil.OID2UUID(oid)
	uuidoTS := time.Unix(uuidOut.Time().UnixTime()).UTC()

	//Print stuff
	fmt.Printf("UUID Out: %s\n", uuidOut.String())
	fmt.Printf("TS Out:\t  %s -> %d\n", uuidoTS.Format(time.RFC3339Nano), uuidoTS.Unix())
}
