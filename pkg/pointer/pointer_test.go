package pointer

import (
	"math"
	"testing"

	"github.com/go-playground/assert/v2"
)

func Test_Value(t *testing.T) {
	var (
		emptyStr         = new(string)
		expectedEmptyStr = ""
	)
	assert.Equal(t, expectedEmptyStr, Value(emptyStr))

	var (
		testStr     = Of("test")
		expectedStr = "test"
	)
	assert.Equal(t, expectedStr, Value(testStr))

	var (
		emptyInt         = new(int)
		expectedEmptyInt = 0
	)
	assert.Equal(t, expectedEmptyInt, Value(emptyInt))

	var (
		testInt     = Of(math.MaxInt32)
		expectedInt = math.MaxInt32
	)
	assert.Equal(t, expectedInt, Value(testInt))

	var (
		emptyArray         = new([26]int)
		expectedEmptyArray = [26]int{}
	)
	assert.Equal(t, expectedEmptyArray, Value(emptyArray))
}
