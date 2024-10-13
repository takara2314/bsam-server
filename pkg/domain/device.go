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

	// マネージャーはデバイス番号が 1 ~ 10 以外でもOK
	if strings.HasPrefix(deviceID, RoleManager) {
		return true
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

func CreateDeviceID(role string, athleteNo int) string {
	return role + strconv.Itoa(athleteNo)
}

func SeparateDeviceIDByRole(deviceIDs []string) ([]string, []string, []string, []string) {
	markIDs := []string{}
	athleteIDs := []string{}
	managerIDs := []string{}
	unknownIDs := []string{}

	for _, deviceID := range deviceIDs {
		role, _, ok := RetrieveRoleAndMarkNo(deviceID)
		if !ok {
			unknownIDs = append(unknownIDs, deviceID)
			continue
		}

		switch role {
		case RoleMark:
			markIDs = append(markIDs, deviceID)
		case RoleAthlete:
			athleteIDs = append(athleteIDs, deviceID)
		case RoleManager:
			managerIDs = append(managerIDs, deviceID)
		}
	}

	return athleteIDs, markIDs, managerIDs, unknownIDs
}
