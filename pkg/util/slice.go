package util

// StrSliceAdd adds new element if this slice does not have it.
func StrSliceAdd(s []string, newStr string) []string {
	for _, str := range s {
		if str == newStr {
			return s
		}
	}

	return append(s, newStr)
}

// StrSliceRemove removes new element if this slice has it.
func StrSliceRemove(s []string, newStr string) []string {
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
