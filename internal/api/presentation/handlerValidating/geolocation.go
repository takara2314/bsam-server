package handlerValidating

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/takara2314/bsam-server/internal/api/presentation/handler"
	"github.com/takara2314/bsam-server/pkg/domain"
)

func GeolocationPOST(c *gin.Context) {
	assocID := c.GetString("assoc_id")

	var req handler.GeolocationPOSTRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "token is required",
		})
		return
	}

	if valid := domain.ValidateDeviceID(req.DeviceID); !valid {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "device_id is invalid",
		})
		return
	}

	handler.GeolocationPOST(c, assocID, req)
}
