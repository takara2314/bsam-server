package handlerValidating

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/takara2314/bsam-server/internal/game/common"
	"github.com/takara2314/bsam-server/pkg/infrastructure/repository/firestore"
	"github.com/takara2314/bsam-server/pkg/racehub"
)

func AssocIDWS(c *gin.Context) {
	var err error
	assocID := c.Param("assocID")

	hub, err := findOrCreateHub(c, assocID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "assoc_id is not found",
		})
		return
	}

	conn, err := racehub.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "failed to upgrade connection",
		})
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
	assoc, err := firestore.FetchAssocByID(ctx, common.FirestoreClient, assocID)
	if err != nil {
		return nil, err
	}

	return racehub.NewHub(assoc.ID), nil
}
