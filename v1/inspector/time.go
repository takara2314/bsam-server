package inspector

import (
	"strings"
	"time"
)

func ParseTimestamp(value string) (time.Time, error) {
	if strings.Contains(value, ".") {
		value = strings.Split(value, ".")[0]
	}

	t, err := time.Parse("2006-01-02 15:04:05 -0700", value+" +0900")

	return t.UTC(), err
}
