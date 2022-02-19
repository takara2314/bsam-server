package utils

import (
	"reflect"
)

// IsNil returns that the value is nil.
func IsNil(value interface{}) bool {
	return reflect.ValueOf(value).String() == "<invalid Value>"
}
