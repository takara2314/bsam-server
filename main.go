package main

import (
	v1 "bsam-server/v1"
	v3 "bsam-server/v3"
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
			"https://bsam-manage.vercel.app",
		}
		router.Use(cors.New(corsConfig))
	} else {
		router.Use(cors.Default())
	}

	v1.Register(router)
	v3.Register(router)

	router.Run(":" + os.Getenv("PORT"))
}
