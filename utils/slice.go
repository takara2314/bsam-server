package utils

// StrSliceAdd adds new element if this slice does not have it.
func StrSliceAdd(s []string, new string) []string {
	for _, str := range s {
		if str == new {
			return s
		}
	}
	return append(s, new)
}

// StrSliceRemove removes new element if this slice has it.
func StrSliceRemove(s []string, new string) []string {
	no := -1

	for i, str := range s {
		if str == new {
			no = i
			break
		}
	}

	if no == -1 {
		return s
	}

	return append(s[:no], s[no+1:]...)
}
