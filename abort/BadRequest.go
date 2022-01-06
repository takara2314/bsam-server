package abort

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BadRequest sends JSON with its message as BadRequest
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, JSON{
		Status:  "BadRequest",
		Message: message,
	})
}
