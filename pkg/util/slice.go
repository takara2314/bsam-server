package util

// AddStrToSliceIfNotExists adds new element if this slice does not have it.
func AddStrToSliceIfNotExists(s []string, newStr string) []string {
	for _, str := range s {
		if str == newStr {
			return s
		}
	}

	return append(s, newStr)
}

// RemoveStrFromSliceIfExists removes new element if this slice has it.
func RemoveStrFromSliceIfExists(s []string, newStr string) []string {
	no := -1

	for i, str := range s {
		if str == newStr {
			no = i
			break
		}
	}

	if no == -1 {
		return s
	}

	return append(s[:no], s[no+1:]...)
}
