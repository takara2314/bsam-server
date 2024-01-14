package utils

import "errors"

// StringSliceToString converts string slice to one string
//
//	{"a", "b", "c"} -> "a, b, c"
func StringSliceToString(s []string) string {
	str := ""
	s_len := len(s)

	for i, item := range s {
		str += item

		if s_len-1 != i {
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
	a_len := len(a)

	if a_len != len(b) {
		return "", errors.New("slice a's length and slice b's length are not the same")
	}

	for i := 0; i < a_len; i++ {
		str += a[i] + " = " + b[i]
		if a_len-1 != i {
			str += ", "
		}
	}

	return str, nil
}
