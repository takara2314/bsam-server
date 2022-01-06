package abort

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Forbidden sends JSON with its message as Forbidden.
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, JSON{
		Status:  "Forbidden",
		Message: message,
	})
}
