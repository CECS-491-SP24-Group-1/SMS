package try

import (
	"errors"
	"fmt"
)

/*
Defines a slim "try" operator, akin to those found in other languages. This
function relies on an inner function calling `panic()` when an error condition
occurs. When doing this, a Go error is rethrown and returned, allowing the
grouping of a myriad of errors into one error object. This function is based
off the implementation found in github.com/ez4o/go-try, but without the
custom "Exception" type.
*/
func Try[T any](f func() T) (val T, err error) {
	//Run the target function
	func() {
		//Catch any panics that may occur and rethrow it as an `error`
		defer func() {
			if r := recover(); r != nil {
				//Construct the error based on the datatype of the recovered panic var
				switch x := r.(type) {
				case string:
					err = errors.New(x)
				case error:
					err = x
				default:
					err = fmt.Errorf("%v", x)
				}
			}
		}()

		//Do something in the function
		val = f()
	}()
	return
}
