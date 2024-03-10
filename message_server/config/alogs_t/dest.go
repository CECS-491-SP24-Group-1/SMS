package alogs_t

import (
	"fmt"
	"strings"
)

// Represents a logging destination for access logs.
type Dest uint8

const (
	dest_parse_err Dest = iota
	SYSOUT
	SYSOUT_FILE
	SYSOUT_RFILE
	SYSERR
	SYSERR_FILE
	SYSERR_RFILE
	FILE
	RFILE
)

var (
	dest_names = map[uint8]string{
		1: "SYSOUT",
		2: "SYSOUT+FILE",
		3: "SYSOUT+RFILE",
		4: "SYSERR",
		5: "SYSERR+FILE",
		6: "SYSERR+RFILE",
		7: "FILE",
		8: "RFILE",
	}
	dest_values = map[string]uint8{
		"SYSOUT":       1,
		"SYSOUT+FILE":  2,
		"SYSOUT+RFILE": 3,
		"SYSERR":       4,
		"SYSERR+FILE":  5,
		"SYSERR+RFILE": 6,
		"FILE":         7,
		"RFILE":        8,
	}
)

// Gets the string equivalent of the object.
func (dest Dest) String() string {
	if dest >= SYSOUT && dest <= RFILE {
		return dest_names[uint8(dest)]
	} else {
		return "<UNKNOWN>"
	}
}

// Converts a string to an access log dest, returns an error if the string is unknown.
func ParseDest(s string) (Dest, error) {
	s = strings.TrimSpace(strings.ToUpper(s))
	value, ok := dest_values[s]
	if !ok {
		return Dest(0), fmt.Errorf("%q is not a valid access log destination", s)
	}
	return Dest(value), nil
}

// Overload for go-toml serialization.
func (dest Dest) MarshalText() ([]byte, error) {
	return []byte(dest.String()), nil
}

// Overload for go-toml deserialization.
func (dest *Dest) UnmarshalText(data []byte) error {
	adest, err := ParseDest(string(data))
	if adest != dest_parse_err {
		*dest = adest
	}
	return err
}
