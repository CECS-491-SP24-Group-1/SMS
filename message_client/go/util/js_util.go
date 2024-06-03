package util

import (
	"syscall/js"

	"github.com/norunners/vert"
)

// Creates a Golang array from a JS array.
func JSArray2GoArray[T any](jsa js.Value, maxlen int, convFunc func(v js.Value) T) []T {
	//Allocate the output array
	arr := make([]T, maxlen)

	//Loop over the JS array and copy the values
	for i := 0; i < min(jsa.Length(), maxlen); i++ {
		arr[i] = convFunc(jsa.Index(i)) //Dynamically calling the correct conversion function
	}

	//Return the array
	return arr
}

// Creates a Golang byte array from a JS array.
func JSArray2GoByteArray(jsa js.Value, maxlen int) []byte {
	return JSArray2GoArray[byte](jsa, maxlen, Val2Any)
}

// Converts a JS object to a Go object using Vert.
func Val2Any[T any](val js.Value) T {
	v := vert.ValueOf(val)
	var out T
	v.AssignTo(&out)
	return out
}

// Creates a generic array from a given array.
func GenerifyArray[T any](arr []T) []interface{} {
	out := make([]interface{}, len(arr))
	for i := 0; i < len(arr); i++ {
		out[i] = interface{}(arr[i])
	}
	return out
}
