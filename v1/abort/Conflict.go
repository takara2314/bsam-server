package abort

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Conflict sends JSON with its message as Conflict.
func Conflict(c *gin.Context, message string) {
	c.JSON(http.StatusConflict, JSON{
		Status:  "Conflict",
		Message: message,
	})
}
