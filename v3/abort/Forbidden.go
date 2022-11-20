package abort

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Forbidden(c *gin.Context) {
	c.JSON(http.StatusForbidden, make(map[string]any))
	c.Abort()
}
