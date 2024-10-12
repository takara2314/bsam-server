package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/takara2314/bsam-server/internal/api/common"
	"github.com/takara2314/bsam-server/pkg/racelib"
	"google.golang.org/grpc/codes"
)

func RacesAssociationIDGET(c *gin.Context) {
	associationID := c.Param("associationID")

	race, code := racelib.FetchLatestRaceDetailByAssociationID(
		c,
		common.FirestoreClient,
		associationID,
	)
	switch code {
	case codes.OK:
		break
	case codes.NotFound:
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "this race is not found",
		})
		return
	case codes.Internal:
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch race",
		})
		return
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch race",
		})
		return
	}

	c.JSON(http.StatusOK, race)
}
