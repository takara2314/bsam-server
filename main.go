package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	route := gin.Default()

	route.GET("/", routeGetFunc)

	route.Run(":" + os.Getenv("PORT"))
}

func routeGetFunc(c *gin.Context) {
	c.String(200, "こんにちは！")
}
