package utils

import (
	"reflect"
)

// IsNil returns that the value is nil.
func IsNil(value interface{}) bool {
	return reflect.ValueOf(value).String() == "<invalid Value>"
}

// StrSliceToAnySlice converts string slice to interface{} slice.
func StrSliceToAnySlice(s []string) []interface{} {
	output := make([]interface{}, len(s))

	for i, str := range s {
		output[i] = str
	}

	return output
}
