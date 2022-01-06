package inspector

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

var (
	permissions map[string]interface{}
)

func init() {
	file, err := ioutil.ReadFile("./permissions.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(file, &permissions)
	if err != nil {
		panic(err)
	}
}
