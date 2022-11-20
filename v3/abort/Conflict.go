package abort

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Conflict(c *gin.Context) {
	c.JSON(http.StatusConflict, make(map[string]any))
	c.Abort()
}
