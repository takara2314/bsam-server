package abort

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Unauthorized sends JSON with its message as Unauthorized.
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, JSON{
		Status:  "Unauthorized",
		Message: message,
	})
}
