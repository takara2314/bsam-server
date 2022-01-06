package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// InfoGET is /user/:username GET request handler
func InfoGET(c *gin.Context) {
	username := c.Param("username")
	c.String(http.StatusOK, "Hello "+username)
}
