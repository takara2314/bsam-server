package inspector

func (ins *Inspector) IsJSON() bool {
	return ins.Request.Header.Get("Content-Type") != "application/json"
}
