package main

import (
	"os"
	"sailing-assist-mie-api/device"
	"sailing-assist-mie-api/user"

	"github.com/gin-gonic/gin"
)

func main() {
	route := gin.Default()

	route.POST("/devices", device.DevicesPOST)

	user.Register(route.Group("/user"))
	route.POST("/users", user.UsersPOST)

	route.Run(":" + os.Getenv("PORT"))
}
