package abort

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BadRequest(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusBadRequest, make(gin.H))
}

func InvalidJSON(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"message": "This request json is invalid.",
	})
}
