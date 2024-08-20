package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthzGET(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}
