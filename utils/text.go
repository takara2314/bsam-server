package utils

import "errors"

var ErrNotSameLengthSlice = errors.New("slice a's length and slice b's length are not the same")

// StringSliceToString converts string slice to one string
//
//	{"a", "b", "c"} -> "a, b, c"
func StringSliceToString(s []string) string {
	str := ""
	sLength := len(s)

	for i, item := range s {
		str += item

		if sLength-1 != i {
			str += ", "
		}
	}

	return str
}

// CreateStrSliceEqualStrSlice creates the string written A element equal B element.
//
//	a = {"a", "b", "c"}
//	b = {"d", "e", "f"}
//	-> "a = d, b = e, c = f"
func CreateStrSliceEqualStrSlice(a []string, b []string) (string, error) {
	str := ""
	aLength := len(a)
	bLength := len(b)

	if aLength != bLength {
		return "", ErrNotSameLengthSlice
	}

	for i := range aLength {
		str += a[i] + " = " + b[i]
		if aLength-1 != i {
			str += ", "
		}
	}

	return str, nil
}
