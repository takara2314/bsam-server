package action

import (
	"log/slog"

	"github.com/bytedance/sonic"
	"github.com/takara2314/bsam-server/pkg/racehub"
)

type RaceAction struct {
	racehub.UnimplementedAction
}

func (r *RaceAction) AuthResult(
	c *racehub.Client,
	output *racehub.AuthResultOutput,
) {
	payload, err := sonic.Marshal(output)
	if err != nil {
		slog.Error(
			"failed to marshal auth result",
			"client", c,
			"error", err,
			"output", output,
		)
		return
	}

	c.Send <- payload

	slog.Info(
		"sent auth result",
		"client", c,
		"output", output,
	)
}
