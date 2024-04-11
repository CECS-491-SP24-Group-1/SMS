package tests

import (
	"encoding/binary"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUUIDxD128(t *testing.T) {
	//Create a new UUID
	uuidIn, _ := uuid.NewV7()
	uuidInStr := strings.ReplaceAll(uuidIn.String(), "-", "")
	fmt.Printf("UUID In:  %s\n", uuidInStr)

	//Create a Decimal128 from the UUID
	d128 := uuid2dec128(uuidIn)
	d128hi, d128lo := d128.GetBytes()
	d128Str := fmt.Sprintf("%016x%016x", d128hi, d128lo)
	fmt.Printf("D128:     %s\n", d128Str)

	//Convert the Decimal128 back to a UUID; the input and output UUID should match up
	uuidOut := dec1282uuid(d128)
	uuidOutStr := strings.ReplaceAll(uuidOut.String(), "-", "")
	fmt.Printf("UUID Out: %s\n", uuidOutStr)

	//Perform equality tests
	if d128Str != uuidOutStr {
		t.Fatalf("string equality: D128 and UUID do not match; %s::%s", d128Str, uuidOutStr)
	}
	if uuidIn != uuidOut {
		t.Fatalf("uuid obj equality: in UUID and out UUID do not match; %s::%s", uuidIn, uuidOut)
	}
}

func uuid2dec128(uuid uuid.UUID) primitive.Decimal128 {
	//Get the high and low values of the UUID
	uuidLo := uuid[0:8]
	uuidHi := uuid[8:16]

	//Convert the byte arrays to uints
	uintLo := binary.BigEndian.Uint64(uuidHi)
	uintHi := binary.BigEndian.Uint64(uuidLo)

	//Create a Decimal128 from the two source uint64s and return
	return primitive.NewDecimal128(uintHi, uintLo)
}

func dec1282uuid(d128 primitive.Decimal128) uuid.UUID {
	//Extract the high and low bytes from the Decimal128
	hi, lo := d128.GetBytes()

	//Create a buffer to hold the two uint64s and copy them in
	buf := make([]byte, 16)
	binary.BigEndian.PutUint64(buf[0:8], hi)
	binary.BigEndian.PutUint64(buf[8:16], lo)

	//Create a UUID from the buffer and return it
	return uuid.UUID(buf)
}
