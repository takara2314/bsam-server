package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func IndexGET(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "B-SAM API Server",
	})
}
