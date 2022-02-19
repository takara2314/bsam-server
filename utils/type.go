package utils

import (
	"reflect"
)

func IsNil(value interface{}) bool {
	return reflect.ValueOf(value).String() == "<invalid Value>"
}
