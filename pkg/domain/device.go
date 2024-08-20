package domain

import (
	"strconv"
	"strings"

	"github.com/takara2314/bsam-server/pkg/util"
)

const (
	RoleMark    = "mark"
	RoleAthlete = "athlete"
	RoleManager = "manager"
	RoleUnknown = "unknown"
)

var idPrefixes = []string{
	RoleMark,
	RoleAthlete,
	RoleManager,
}

func ValidateDeviceID(deviceID string) bool {
	if !util.HasAnyPrefix(deviceID, idPrefixes) {
		return false
	}

	iotaInRole, ok := strconv.Atoi(util.StripAnyPrefix(deviceID, idPrefixes))
	if ok != nil {
		return false
	}

	if iotaInRole <= 0 || iotaInRole > 10 {
		return false
	}

	return true
}

func RetrieveRoleAndMarkNo(deviceID string) (string, int, bool) {
	if !ValidateDeviceID(deviceID) {
		return "", -1, false
	}

	for _, prefix := range idPrefixes {
		if strings.HasPrefix(deviceID, prefix) {
			MarkNo := -1

			if prefix == RoleMark {
				var err error
				MarkNo, err = strconv.Atoi(
					strings.TrimPrefix(deviceID, prefix),
				)
				if err != nil {
					return "", 0, false
				}
			}

			return prefix, MarkNo, true
		}
	}

	return "", -1, false
}

func CreateDeviceID(role string, markNo int) string {
	return role + strconv.Itoa(markNo)
}
