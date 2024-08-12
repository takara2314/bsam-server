package handlerValidating

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/takara2314/bsam-server/internal/game/common"
	"github.com/takara2314/bsam-server/internal/game/presentation/action"
	"github.com/takara2314/bsam-server/internal/game/presentation/event"
	"github.com/takara2314/bsam-server/internal/game/presentation/handler"
	repoFirestore "github.com/takara2314/bsam-server/pkg/infrastructure/repository/firestore"
	"github.com/takara2314/bsam-server/pkg/racehub"
	"github.com/takara2314/bsam-server/pkg/taskmanager"
)

func AssociationIDWS(c *gin.Context) {
	var err error
	associationID := c.Param("associationID")

	hub, err := findOrCreateHub(c, associationID)
	if err != nil {
		slog.Warn(
			"failed to find or create hub",
			"association_id", associationID,
			"error", err,
		)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "association_id is not found",
		})
		return
	}

	conn, err := racehub.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Warn(
			"failed to upgrade connection",
			"association_id", associationID,
			"error", err,
		)
		return
	}

	hub.Register(conn)
}

func findOrCreateHub(c *gin.Context, associationID string) (*racehub.Hub, error) {
	if hub, exist := common.Hubs[associationID]; exist {
		return hub, nil
	}

	hub, err := createNewHub(c, associationID)
	if err != nil {
		return nil, err
	}

	common.Hubs[associationID] = hub
	return hub, nil
}

func createNewHub(ctx context.Context, associationID string) (*racehub.Hub, error) {
	association, err := repoFirestore.FetchAssociationByID(
		ctx,
		common.FirestoreClient,
		associationID,
	)
	if err != nil {
		return nil, err
	}

	tm := taskmanager.NewManager(common.FirestoreClient)

	return racehub.NewHub(
		association.ID,
		tm,
		&event.RaceEvent{},
		&handler.RaceHandler{},
		&action.RaceAction{},
	), nil
}
