package util

import (
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

// Checks if an integer is inside of a specified range.
func InRange(num int64, min int64, max int64) bool {
	return num >= min && num <= max
}

/*
Generates a random string of size n, given a character set. By default,
this will be all alphanumeric characters.
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
