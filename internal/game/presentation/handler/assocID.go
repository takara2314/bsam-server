package handler

import (
	"fmt"

	"github.com/takara2314/bsam-server/pkg/racehub"
)

type RaceHandler struct {
	racehub.UnimplementedHandler
}

func (r *RaceHandler) Auth(c *racehub.Client) {
	fmt.Println("auth")
}

func (r *RaceHandler) PostGeolocation(c *racehub.Client) {
	fmt.Println("post geolocation")
}
