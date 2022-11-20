package abort

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, make(map[string]any))
	c.Abort()
}
