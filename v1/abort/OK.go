package abort

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// OK sends JSON with its message as OK.
func OK(c *gin.Context, message string) {
	c.JSON(http.StatusOK, JSON{
		Status:  "OK",
		Message: message,
	})
}
