package tests

import (
	"testing"
)

func testType[T any](t *testing.T, targ T, expected bool) {
	actual := isComplexType(targ)
	if expected != actual {
		t.Errorf("Unexpected test output %v for item `%v`; expected %v\n", expected, targ, actual)
	}
}

// See the following for the list of basic types https://go.dev/tour/basics/11
func TestType(t *testing.T) {
	//Integral types return false
	testType(t, false, false)
	testType(t, byte(37), false)
	testType(t, int(373), false)
	testType(t, int8(33), false)
	testType(t, int16(3336), false)
	testType(t, int32(352336), false)
	testType(t, int64(35226336), false)
	testType(t, uint(373), false)
	testType(t, uint8(33), false)
	testType(t, uint16(3336), false)
	testType(t, uint32(352336), false)
	testType(t, uint64(35226336), false)
	testType(t, uintptr(35226336), false)

	//Float and complex types return false
	testType(t, float32(262.2156), false)
	testType(t, float64(262.2156), false)
	testType(t, complex64(262.2156), false)
	testType(t, complex128(262.2156), false)

	//Characters and strings return false
	testType(t, '5', false)
	testType(t, "hello", false)

	//Arrays and slices of the above should return false
	array1 := [2]string{"hello", "world"}
	array2 := [3]int{1, 2, 3}
	slice1 := make([]float32, 1)
	slice1[0] = 3.14
	slice2 := make([]rune, 1)
	slice2[0] = 'x'
	testType(t, array1, false)
	testType(t, array2, false)
	testType(t, slice1, false)
	testType(t, slice2, false)
}
