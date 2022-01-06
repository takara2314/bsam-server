package inspector

import (
	"strings"
)

// isExistPermission check that its permission is exist
func isExistPermission(permission string) bool {
	splitted := strings.Split(permission, ".")

	var nest map[string]interface{} = permissions

	for index, item := range splitted {
		if item == "*" {
			return true
		}

		if _, exist := nest[item]; !exist {
			return false
		}

		// Illegal that $nest's type is string even if $index is not final loop
		switch nest[item].(type) {
		case map[string]interface{}:
			nest = nest[item].(map[string]interface{})
		case string:
			if index != len(splitted)-1 {
				return false
			}
		}
	}

	return true
}

// HasPermission check that its user have this permission
func (ins Inspector) HasPermission(permission string) bool {
	if !isExistPermission(permission) {
		return false
	}

	return true
}
