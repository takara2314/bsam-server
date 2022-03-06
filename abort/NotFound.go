package abort

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NotFound sends JSON with its message as NotFound.
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, JSON{
		Status:  "Not Found",
		Message: message,
	})
}
