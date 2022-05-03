package main

import (
	"os"
	"sailing-assist-mie-api/device"
	"sailing-assist-mie-api/group"
	"sailing-assist-mie-api/race"
	"sailing-assist-mie-api/user"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// CORS settings
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000", "http://sailing-assist-mie-manage.herokuapp.com", "https://sailing-assist-mie-manage.herokuapp.com"}
	router.Use(cors.New(corsConfig))

	// Device API
	device.Register(router.Group("/device"))

	// User API
	user.Register(router.Group("/user"))
	router.POST("/users", user.UsersPOST)

	// Race API and Race Socket
	race.Register(router.Group("/race"))
	router.GET("/races", race.RacesGET)
	router.POST("/races", race.RacesPOST)
	router.GET("/racing/:id", race.RacingWS)
	go race.AutoRooming()

	// Group API
	group.Register(router.Group("/group"))
	router.POST("/groups", group.GroupsPOST)

	router.Run(":" + os.Getenv("PORT"))
}
