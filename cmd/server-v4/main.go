package main

import (
	"log"
	"os"

	server "github.com/takara2314/bsam-server/v4"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// CORS settings
	if os.Getenv("GIN_MODE") == gin.ReleaseMode {
		corsConfig := cors.DefaultConfig()
		corsConfig.AllowOrigins = []string{
			os.Getenv("RACE_MONITOR_SITE_URL"),
			os.Getenv("TEST_SITE_URL"),
		}
		router.Use(cors.New(corsConfig))
	} else {
		router.Use(cors.Default())
	}

	server.Register(router)

	log.Printf("Server is running on port %s\n", os.Getenv("PORT"))

	err := router.Run(":" + os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}
}
