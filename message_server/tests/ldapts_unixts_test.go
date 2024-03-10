package tests

import (
	"fmt"
	"testing"
	"time"
)

func TestLDAPTSxUnixTS(t *testing.T) {
	//Get the current time as a Unix timestamp
	unixIn := time.Now().Unix()

	//Convert to LDAP and back to Unix
	ldap := unix2ldap(unixIn)
	unixOut := ldap2unix(ldap)

	//Print results
	fmt.Printf("Unix TS I: %d\n", unixIn)
	fmt.Printf("LDAP TS:   %d\n", ldap)
	fmt.Printf("Unix TS O: %d\n", unixOut)

	//Test for equality
	if unixIn != unixOut {
		t.Errorf("mismatched Unix timestamp; %d::%d", unixIn, unixOut)
		t.FailNow()
	}
}

// See: https://gist.github.com/caseydunham/508e2994e1195e4cb8e4
func unix2ldap(uepoch int64) int64 {
	//Calculate the number of seconds between the LDAP epoch and the Unix epoch
	ADToUnixConverter := ((1970-1601)*365 - 3 + round((1970-1601)/4)) * 86400

	//Calculate the number of seconds after the LDAP epoch
	secsAfterADEpoch := ADToUnixConverter + uepoch

	//Add on the nanoseconds
	return secsAfterADEpoch * 1e7
}

// See: https://gist.github.com/caseydunham/508e2994e1195e4cb8e4
func ldap2unix(lepoch int64) int64 {
	//Remove the nanoseconds
	secsAfterADEpoch := lepoch / 1e7

	//Calculate the number of seconds between the LDAP epoch and the Unix epoch
	ADToUnixConverter := ((1970-1601)*365 - 3 + round((1970-1601)/4)) * 86400

	//Calculate the number of seconds after the LDAP epoch
	return int64(secsAfterADEpoch - ADToUnixConverter)
}

func round(num float64) int64 {
	if num < 0 {
		return int64(num - 0.5)
	}
	return int64(num + 0.5)
}
