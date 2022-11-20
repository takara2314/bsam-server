package main

import (
	v3 "bsam-server/v3"
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
			os.Getenv("MANAGE_SITE"),
		}
		router.Use(cors.New(corsConfig))
	} else {
		router.Use(cors.Default())
	}

	v3.Register(router)

	fmt.Printf("Server is running on port %s\n", os.Getenv("PORT"))
	router.Run(":" + os.Getenv("PORT"))
}
