package ip_addr

import (
	"fmt"
	"net"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

const (
	//The size of an IPv4 address in bytes.
	IP4_SIZE = 4

	//The size of an IPv6 in bytes.
	IP6_SIZE = 16
)

//
//-- ENUM: IPType
//

// Represents the type of an IP address.
type IPType uint8

const (
	IP_TYPE4 IPType = iota
	IP_TYPE6
)

//
//-- Class: IPAddr
//

/*
Represents an IP address. Aims to be interoperable with the std's `net.IP`
while also being a fixed size, unlike `net.IP` which is not a fixed size.
*/
type IPAddr struct {
	//The raw bytes of the IP address.
	Bytes [IP6_SIZE]byte `json:"bytes" bson:"bytes"`

	//The type of IP address this is.
	Type IPType `json:"type" bson:"type"`
}

// Converts a `net.IP` object to an `IPAddr` object.
func FromNetIP(ip net.IP) IPAddr {
	//Create the object and return
	obj := IPAddr{
		Type: TypeOf(ip),
	}
	copy(obj.Bytes[:], []byte(ip)[:IP6_SIZE])
	return obj
}

// Converts a string IP to an `IPAddr` object.
func FromString(ip string) IPAddr {
	//Derive a net.IP object
	nip := net.ParseIP(ip)

	//Create the object and return
	obj := IPAddr{
		Type: TypeOf(nip),
	}
	copy(obj.Bytes[:], []byte(nip)[:IP6_SIZE])
	return obj
}

// Converts an `IPAddr` object to a string.
func (ip IPAddr) String() string {
	return ip.ToNetIP().String()
}

// Converts an `IPAddr` object to a byte array (17 bytes).
func (ip IPAddr) ToBytes() []byte {
	bytes := [IP6_SIZE + 1]byte{}
	copy(bytes[:], ip.Bytes[0:16])
	bytes[16] = byte(ip.Type)
	return bytes[:]
}

// Converts a `net.IP` object to an `IPAddr` object.
func (ip IPAddr) ToNetIP() net.IP {
	return net.IP(ip.Bytes[:])
}

// MarshalBSONValue implements the bson.ValueMarshaler interface.
func (ip IPAddr) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.TypeBinary, bsoncore.AppendBinary(nil, bson.TypeBinaryGeneric, ip.ToBytes()), nil
}

// Marshals an IP to text. Used downstream by JSON and BSON marshalling.
func (ip IPAddr) MarshalText() (text []byte, err error) {
	return []byte(ip.String()), nil
}

// UnmarshalBSONValue implements the bson.ValueUnmarshaler interface.
func (ip *IPAddr) UnmarshalBSONValue(t bsontype.Type, raw []byte) error {
	//Ensure the incoming type is correct
	if t != bson.TypeBinary {
		return fmt.Errorf("(IPAddr) invalid format on unmarshalled bson value")
	}

	//Read the data from the BSON item
	_, data, _, ok := bsoncore.ReadBinary(raw)
	if !ok {
		return fmt.Errorf("(IPAddr) not enough bytes to unmarshal bson value")
	}
	copy(ip.Bytes[:], data)

	//No errors, so return nil
	return nil
}

// Unmarshals an IP from a string. Used downstream by JSON and BSON marshalling.
func (ip *IPAddr) UnmarshalText(text []byte) error {
	*ip = FromString(string(text)[:])
	return nil
}
