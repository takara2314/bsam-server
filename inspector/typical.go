package inspector

func (ins *Inspector) IsJSON() bool {
	if ins.Request.Header.Get("Content-Type") != "application/json" {
		return false
	}

	return true
}
