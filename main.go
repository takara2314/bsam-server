package main

import (
	v4 "bsam-server/v4"
	"fmt"
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
			os.Getenv("RACE_MONITOR_SITE_URL"),
			os.Getenv("TEST_SITE_URL"),
		}
		router.Use(cors.New(corsConfig))
	} else {
		router.Use(cors.Default())
	}

	v4.Register(router)

	fmt.Printf("Server is running on port %s\n", os.Getenv("PORT"))
	router.Run(":" + os.Getenv("PORT"))
}
