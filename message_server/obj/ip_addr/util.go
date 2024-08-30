package ip_addr

import (
	"net"
	"strings"
)

// Converts a string IP from `http.Request` to an `IPAddr` object.
func HttpIP2IPAddr(ip string) IPAddr {
	return FromNetIP(HttpIP2NetIP(ip))
}

// Converts a string IP from `http.Request` to a `net.IP` object.
func HttpIP2NetIP(ip string) net.IP {
	rawIP := ip[0:strings.LastIndex(ip, ":")] //Get just the IP; last colon indicates the port
	return net.ParseIP(rawIP)
}

// Determines the type of a `net.IP` address.
func TypeOf(ip net.IP) IPType {
	/*
		IPv6 addresses cannot be represented as IPv4 strings, so `net.To16()`
		returns `nil` if the IP is an IPv4. We are exploiting this behavior.
	*/
	if ip.To4() != nil {
		return IP_TYPE4
	} else if ip.To16() != nil {
		return IP_TYPE6
	}
	return IP_TYPE6
}
