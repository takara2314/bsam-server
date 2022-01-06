package main

import (
	"os"

	"sailing-assist-mie-api/user"

	"github.com/gin-gonic/gin"
)

func main() {
	route := gin.Default()

	user.Register(route.Group("/user"))
	route.POST("/users", user.UsersPOST)

	route.Run(":" + os.Getenv("PORT"))
}
