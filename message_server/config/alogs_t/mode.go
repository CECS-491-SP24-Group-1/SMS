package alogs_t

import (
	"fmt"
	"strings"
)

// Represents a logging mode for access logs.
type Mode uint8

const (
	mode_parse_err Mode = iota
	OFF
	FMT
	FMT_SIM
	JSON
)

var (
	mode_names = map[uint8]string{
		1: "OFF",
		2: "FMT",
		3: "FMT_SIM",
		4: "JSON",
	}
	mode_values = map[string]uint8{
		"OFF":     1,
		"FMT":     2,
		"FMT_SIM": 3,
		"JSON":    4,
	}
)

// Gets the string equivalent of the object.
func (alm Mode) String() string {
	if alm >= OFF && alm <= JSON {
		return mode_names[uint8(alm)]
	} else {
		return "<UNKNOWN>"
	}
}

// Converts a string to an access log dest, returns an error if the string is unknown.
func ParseALM(s string) (Mode, error) {
	s = strings.TrimSpace(strings.ToUpper(s))
	value, ok := mode_values[s]
	if !ok {
		return Mode(0), fmt.Errorf("%q is not a valid access log mode", s)
	}
	return Mode(value), nil
}

// Overload for go-toml serialization.
func (alm Mode) MarshalText() ([]byte, error) {
	return []byte(alm.String()), nil
}

// Overload for go-toml deserialization.
func (alm *Mode) UnmarshalText(data []byte) error {
	amode, err := ParseALM(string(data))
	if amode != mode_parse_err {
		*alm = amode
	}
	return err
}
