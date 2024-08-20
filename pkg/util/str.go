package util

import "strings"

func HasAnyPrefix(str string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(str, prefix) {
			return true
		}
	}
	return false
}

func StripAnyPrefix(str string, prefixes []string) string {
	for _, prefix := range prefixes {
		if strings.HasPrefix(str, prefix) {
			return strings.TrimPrefix(str, prefix)
		}
	}
	return str
}

func FindPrefixIfHasAnyPrefix(str string, prefixes []string) string {
	for _, prefix := range prefixes {
		if strings.HasPrefix(str, prefix) {
			return prefix
		}
	}
	return ""
}
