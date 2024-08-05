package domain

import (
	"strconv"

	"github.com/takara2314/bsam-server/pkg/util"
)

var idPrefixes = []string{
	"mark",
	"athlete",
	"manager",
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
