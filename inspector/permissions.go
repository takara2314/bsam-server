package inspector

import (
	"sailing-assist-mie-api/utils"
	"strings"
)

// isExistPermission check that its permission is exist.
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

		// Illegal that $nest's type is string even if $index is not final loop.
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

// HasPermission check that its token has this permission.
// Admit if its token has any of them.
func (ins *Inspector) HasPermission(permissions []string) bool {
	isFound := false

	for _, permission := range permissions {
		if isExistPermission(permission) {
			isFound = true
			break
		}
	}

	if !isFound || !ins.IsTokenEnabled {
		return false
	}

	for _, has := range ins.Token.Permissions {
		for _, req := range permissions {
			hasSplitted := strings.Split(has, ".")
			reqSplitted := strings.Split(req, ".")

			loop := utils.MaxInt(len(hasSplitted), len(reqSplitted))
			for i := 0; i < loop; i++ {
				// Containing "*" is special permission.
				if hasSplitted[i] == "*" {
					return true
				}

				if hasSplitted[i] != reqSplitted[i] {
					break
				}
			}
		}
	}

	return false
}
