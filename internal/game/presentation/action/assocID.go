package action

import (
	"github.com/takara2314/bsam-server/pkg/racehub"
)

type RaceAction struct {
	racehub.UnimplementedAction
}

func (r *RaceAction) AuthResult(
	c *racehub.Client,
	ok bool,
	message string,
) (*racehub.AuthResultOutput, error) {
	return &racehub.AuthResultOutput{
		MessageType: racehub.ActionTypeAuthResult,
		OK:          ok,
		DeviceID:    c.DeviceID,
		Role:        c.Role,
		MarkNo:      c.MarkNo,
		Message:     message,
	}, nil
}

func (r *RaceAction) MarkGeolocations(
	c *racehub.Client,
) (*racehub.MarkGeolocationsOutput, error) {
	panic("not implemented")
}
