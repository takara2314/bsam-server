package utils

import "errors"
import "strings"

var ErrNotSameLengthSlice = errors.New("slice a's length and slice b's length are not the same")

// StringSliceToString converts string slice to one string
//
//	{"a", "b", "c"} -> "a, b, c"
func StringSliceToString(s []string) string {
	return strings.Join(s, ", ")
}

// CreateStrSliceEqualStrSlice creates the string written A element equal B element.
//
//	a = {"a", "b", "c"}
//	b = {"d", "e", "f"}
//	-> "a = d, b = e, c = f"
func CreateStrSliceEqualStrSlice(a []string, b []string) (string, error) {
	aLength := len(a)
	bLength := len(b)

	if aLength != bLength {
		return "", ErrNotSameLengthSlice
	}

	pairs := make([]string, 0, aLength)
	for i := range aLength {
		pairs = append(pairs, a[i]+" = "+b[i])
	}

	return strings.Join(pairs, ", "), nil
}
