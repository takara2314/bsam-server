package main

import (
	"os"
	"sailing-assist-mie-api/device"
	"sailing-assist-mie-api/race"
	"sailing-assist-mie-api/user"

	"github.com/gin-gonic/gin"
)

func main() {
	route := gin.Default()

	device.Register(route.Group("/device"))

	user.Register(route.Group("/user"))
	route.POST("/users", user.UsersPOST)

	race.Register(route.Group("/race"))
	route.POST("/races", race.RacesPOST)

	route.Run(":" + os.Getenv("PORT"))
}
