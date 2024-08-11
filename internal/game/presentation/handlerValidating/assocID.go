package handlerValidating

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/takara2314/bsam-server/internal/game/common"
	"github.com/takara2314/bsam-server/internal/game/presentation/action"
	"github.com/takara2314/bsam-server/internal/game/presentation/handler"
	repoFirestore "github.com/takara2314/bsam-server/pkg/infrastructure/repository/firestore"
	"github.com/takara2314/bsam-server/pkg/racehub"
)

func AssocIDWS(c *gin.Context) {
	var err error
	assocID := c.Param("assocID")

	hub, err := findOrCreateHub(c, assocID)
	if err != nil {
		slog.Warn(
			"failed to find or create hub",
			"assoc_id", assocID,
			"error", err,
		)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "assoc_id is not found",
		})
		return
	}

	conn, err := racehub.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Warn(
			"failed to upgrade connection",
			"assoc_id", assocID,
			"error", err,
		)
		return
	}

	hub.Register(conn)
}

func findOrCreateHub(c *gin.Context, assocID string) (*racehub.Hub, error) {
	if hub, exist := common.Hubs[assocID]; exist {
		return hub, nil
	}

	hub, err := createNewHub(c, assocID)
	if err != nil {
		return nil, err
	}

	common.Hubs[assocID] = hub
	return hub, nil
}

func createNewHub(ctx context.Context, assocID string) (*racehub.Hub, error) {
	assoc, err := repoFirestore.FetchAssocByID(
		ctx,
		common.FirestoreClient,
		assocID,
	)
	if err != nil {
		return nil, err
	}

	return racehub.NewHub(
		assoc.ID,
		&handler.RaceHandler{},
		&action.RaceAction{},
	), nil
}
