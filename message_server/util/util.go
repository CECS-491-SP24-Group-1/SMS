package util

import (
	crand "crypto/rand"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
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
Swallows an `error` return on a function, running `panic()` if one occurs.
This function should be used on functions that are known by the programmer
to not return errors on runtime, such as for `regexp.Compile()` on hardcoded
regexps. Adapted from: https://stackoverflow.com/a/73584801
*/
func Must[T any](obj T, err error) T {
	if err != nil {
		panic(err)
	}
	return obj
}

/*
Swallows an `error` return on a function, running `panic()` if one occurs.
This function should be used on functions that are known by the programmer
to not return errors on runtime, such as for `regexp.Compile()` on hardcoded
regexps. Adapted from: https://stackoverflow.com/a/73584801
*/
func MustUnit(err error) {
	if err != nil {
		panic(err)
	}
}

/*
GenerateRandomBytes returns securely generated random bytes. This function
is an alias of `util.Must(util.GenRandBytes(int))`.
*/
func MustGenRandBytes(n int) []byte {
	return Must(GenRandBytes(n))
}

// Shorthand for `time.Now().Truncate(time.Millisecond).UTC()`.
func NowMillis() time.Time {
	return time.Now().Truncate(time.Millisecond).UTC()
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

/*
Replaces n characters on either side of a string with asterisks,
effectively redacting the contents of the string.
*/
//TODO: redo using string slicing
func Redact(str string, n int) string {
	//Convert the string to a slice of runes
	runes := []rune(str)

	//Replace the nth characters from both ends with asterisks
	for i := 0; i < n; i++ {
		if len(runes) > i {
			runes[i] = '*'
		}
		if len(runes) > len(runes)-i-1 {
			runes[len(runes)-i-1] = '*'
		}
	}

	//Convert the slice of runes back to a string
	return string(runes)
}

/*
Redacts an email address, leaving only one character on either side of both
the email account and the domain name. Example: `johndoe@example.com` becomes
`j*****e@e*****e.com`. See the following StackExchange article for justification
of the format used: https://security.stackexchange.com/a/213700
*/
func RedactEmail(email string) string {
	//Split the email at the last occurrence of the `@` rune
	name, domain := SplitAtLastRune(email, '@')

	//Split the domain at the last occurrence of the `.` rune
	dname, tld := SplitAtLastRune(domain, '.')

	//Redact the email and domain names
	name = RedactCenter(name, 1)
	dname = RedactCenter(dname, 1)

	//Recombine the redacted name and domain
	return name + "@" + dname + "." + tld
}

/*
Replaces n characters from the center of a string with asterisks,
effectively redacting the contents of the string.
*/
func RedactCenter(str string, n int) string {
	if len(str) <= n*2 {
		return strings.Repeat("*", len(str))
	} else {
		return str[0:n] + strings.Repeat("*", len(str)-(n*2)) + str[len(str)-n:]
	}
}

// Splits a string into two pieces at the first position of a given rune.
func SplitAtFirstRune(s string, r rune) (string, string) {
	//Find the first index of the rune
	firstIndex := strings.Index(s, string(r))
	if firstIndex == -1 {
		//Rune not found, return the original string and an empty string
		return s, ""
	}
	//Split the string into two parts at the position of the rune
	return s[:firstIndex], s[firstIndex+1:]
}

// Splits a string into two pieces at the last position of a given rune.
func SplitAtLastRune(s string, r rune) (string, string) {
	//Find the last index of the rune
	lastIndex := strings.LastIndex(s, string(r))
	if lastIndex == -1 {
		//Rune not found, return the original string and an empty string
		return s, ""
	}
	//Split the string into two parts at the position of the rune
	return s[:lastIndex], s[lastIndex+1:]
}

// Truncates a time to milliseconds, chopping off any micro or nano seconds.
func Strip2Millis(t time.Time) time.Time {
	return t.Truncate(time.Millisecond)
}

/*
Calculates the offset from UTC based on output from the
`Date.prototype.getTimezoneOffset()`. JS function.
*/
func Time2Offset(tin time.Time, off int) time.Time {
	/*
		Ensure the offset is in the range +/- 720 since `Date.prototype.getTimezoneOffset()`
		returns the offset from UTC in minutes and 60 * 12 = 720. See the following webpage
		for more info: https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Date/getTimezoneOffset
	*/
	if off > 720 || off < -720 {
		return tin
	}

	//Calculate the offset from UTC of the given offset and get the input time's new value
	loc := time.FixedZone("", off*60)
	return tin.UTC().In(loc)
}

/*
Calculates the offset from UTC based on the value of the `X-Timezone-Offset`
HTTP header.
*/
func Time2OffsetReq(tin time.Time, r *http.Request) time.Time {
	//Ensure the string is not blank or an invalid integer
	off := r.Header.Get("X-Timezone-Offset")
	ioff, err := strconv.Atoi(off)
	if err != nil {
		return tin
	}
	return Time2Offset(tin, ioff)
}
