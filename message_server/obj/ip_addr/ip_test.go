package ip_addr

import (
	"fmt"
	"net"
	"testing"
)

func TestFromNetIP(t *testing.T) {
	//Set the starting IP addresses (Cloudflare DNS)
	i4 := net.ParseIP("1.1.1.1")
	i6 := net.ParseIP("2606:4700:4700::1111")
	//i6 := net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334")

	//Ensure the raw byte arrays are the same
	if [IP6_SIZE]byte(i4) != FromNetIP(i4).Bytes {
		t.Errorf("Unequal i4 arrays: '%v' & '%v'", []byte(i4), FromNetIP(i4).Bytes)
	}
	if [IP6_SIZE]byte(i6) != FromNetIP(i6).Bytes {
		t.Errorf("Unequal i6 arrays: '%v' & '%v'", []byte(i6), FromNetIP(i6).Bytes)
	}
}

func TestIPTypeOf(t *testing.T) {
	//Define tests and expected results
	tests := []string{
		"1.1.1.1", "2606:4700:4700::1111", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", "127.0.0.1",
	}
	expected := []IPType{
		IP_TYPE4, IP_TYPE6, IP_TYPE6, IP_TYPE4,
	}

	//Run tests
	for i := 0; i < len(tests); i++ {
		result := TypeOf(net.ParseIP(tests[i]))
		if result != expected[i] {
			t.Fatalf("Incorrect type received. Got %d, expected %d", result, expected[i])
			t.FailNow()
		}
	}
}

func Test2String(t *testing.T) {
	//Set the starting IP addresses (Cloudflare DNS)
	i4 := net.ParseIP("1.1.1.1")
	i6 := net.ParseIP("2606:4700:4700::1111")
	//i6 := net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334")

	//Print strings
	fmt.Printf("I4: %s\n", FromNetIP(i4))
	fmt.Printf("I6: %s\n", FromNetIP(i6))
}
