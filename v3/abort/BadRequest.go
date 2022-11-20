package abort

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BadRequest(c *gin.Context) {
	c.JSON(http.StatusBadRequest, make(map[string]any))
	c.Abort()
}
