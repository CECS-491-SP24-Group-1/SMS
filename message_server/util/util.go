package util

import (
	crand "crypto/rand"
	"encoding/base64"
	"fmt"
	"math/rand"
)

// Checks if any item in a given list is equal to one singular given item.
func EqualsAny[T comparable](targ T, items ...T) bool {
	for _, item := range items {
		if targ == item {
			return true
		}
	}
	return false
}

/*
Formats a uint denoting a size as a string down to the nearest whole size
unit up to TB. For example, the value `827607531` will become `789.27MB`.
*/
func FormatBytes(size uint64, sp bool) string {
	//Determine the correct space character
	ws := ""
	if sp {
		ws = " "
	}

	//Get the size as a float64
	value := float64(size)

	//Determine the correct units
	unit := ""
	switch {
	case size >= 1<<40:
		unit = "TB"
		value /= 1 << 40
	case size >= 1<<30:
		unit = "GB"
		value /= 1 << 30
	case size >= 1<<20:
		unit = "MB"
		value /= 1 << 20
	case size >= 1<<10:
		unit = "KB"
		value /= 1 << 10
	default:
		return fmt.Sprintf("%d%sB", size, ws)
	}

	//Sprintf the output
	return fmt.Sprintf("%.2f%s%s", value, ws, unit)
}

/*
GenerateRandomBytes returns securely generated random bytes. It will
return an error if the system's secure random number generator fails
to function correctly, in which case the caller should not continue.
Function is sourced from the following website:
https://blog.questionable.services/article/generating-secure-random-numbers-crypto-rand/
*/
func GenRandBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := crand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

/*
GenRandString returns a URL-safe, base64 encoded securely generated
random string. It will return an error if the system's secure random
number generator fails to function correctly, in which case the caller
should not continue. Function is sourced from the following website:
https://blog.questionable.services/article/generating-secure-random-numbers-crypto-rand/
*/
func GenRandString(s int) (string, error) {
	b, err := GenRandBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

// Checks if an integer is inside of a specified range.
func InRange(num int64, min int64, max int64) bool {
	return num >= min && num <= max
}

/*
Generates a random string of size n, given a character set. By default,
this will be all alphanumeric characters. Keep in mind that this function
is NOT cryptographically secure, and should not be used for generating
sensitive data.
*/
func RandomString(size int, charset string) string {
	//Check if the charset is empty
	if charset == "" {
		charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	//Preallocate the output string to save on space
	str := make([]byte, size)

	//Create the random string by choosing a random character from the charset n times
	for i := range str {
		str[i] = charset[rand.Intn(len(charset))]
	}

	//Return the resultant string
	return string(str)
}
