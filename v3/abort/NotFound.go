package abort

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NotFound(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusNotFound, make(gin.H))
}
