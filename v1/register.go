package v1

import (
	"bsam-server/v1/auth"
	"bsam-server/v1/device"
	"bsam-server/v1/group"
	"bsam-server/v1/race"
	"bsam-server/v1/user"

	"github.com/gin-gonic/gin"
)

func Register(e *gin.Engine) *gin.RouterGroup {
	router := e.Group("/")

	// Device API
	device.Register(router.Group("/device"))

	// User API
	user.Register(router.Group("/user"))
	router.GET("/users", user.UsersGET)
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

	// Authorization API
	auth.Register(router.Group("/auth"))

	return router
}
