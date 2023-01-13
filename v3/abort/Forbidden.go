package abort

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Forbidden(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusForbidden, make(gin.H))
}
