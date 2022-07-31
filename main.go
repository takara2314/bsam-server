package main

import (
	v1 "bsam-server/v1"
	v2 "bsam-server/v2"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// CORS settings
	if os.Getenv("GIN_MODE") == "release" {
		corsConfig := cors.DefaultConfig()
		corsConfig.AllowOrigins = []string{
			"http://sailing-assist-mie-manage.herokuapp.com",
			"https://sailing-assist-mie-manage.herokuapp.com",
		}
		router.Use(cors.New(corsConfig))
	} else {
		router.Use(cors.Default())
	}

	v1.Register(router)
	v2.Register(router)

	router.Run(":" + os.Getenv("PORT"))
}
