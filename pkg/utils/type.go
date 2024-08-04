package utils

import (
	"reflect"
)

// IsNil returns that the value is nil.
func IsNil(value any) bool {
	return reflect.ValueOf(value).String() == "<invalid Value>"
}

// StrSliceToAnySlice converts string slice to any slice.
func StrSliceToAnySlice(s []string) []any {
	output := make([]any, len(s))

	for i, str := range s {
		output[i] = str
	}

	return output
}
