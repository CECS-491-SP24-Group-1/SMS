package ip_addr

import "net"

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

// Converts a `net.IP` object to an `IPAddr` object.
func (ip IPAddr) ToNetIP() net.IP {
	return net.IP(ip.Bytes[:])
}
