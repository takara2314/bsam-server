package abort

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Conflict(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusConflict, make(gin.H))
}
