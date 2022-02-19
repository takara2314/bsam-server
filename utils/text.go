package utils

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
