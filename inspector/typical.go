package inspector

import "strings"

func (ins *Inspector) IsJSON() bool {
	return strings.Contains(ins.Request.Header.Get("Content-Type"), "application/json")
}
