package inspector

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

var (
	permissions map[string]interface{}

	ErrTokenNotFound            = errors.New("the token is not found")
	ErrAuthorizationTypeInvalid = errors.New("this authorization type is invalid")
	ErrWrongToken               = errors.New("this token is wrong")
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
