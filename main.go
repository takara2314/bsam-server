package main

import (
	"os"
	"sailing-assist-mie-api/device"
	"sailing-assist-mie-api/group"
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
	route.GET("/races", race.RacesGET)
	route.POST("/races", race.RacesPOST)
	route.GET("/racing/:id", race.RacingWS)
	go race.AutoRooming()

	group.Register(route.Group("/group"))
	route.POST("/groups", group.GroupsPOST)

	route.Run(":" + os.Getenv("PORT"))
}
