package handlerRouting

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

	if _, exist := common.Hubs[assocID]; !exist {
		common.Hubs[assocID], err = createNewHub(c, assocID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "assoc_id is not found",
			})
			return
		}
	}

	conn, err := racehub.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	common.Hubs[assocID].Register(conn)
}

func createNewHub(ctx context.Context, assocID string) (*racehub.Hub, error) {
	assoc, err := firestore.FetchAssocByID(ctx, common.FirestoreClient, assocID)
	if err != nil {
		return nil, err
	}

	return racehub.NewHub(assoc.ID), nil
}
