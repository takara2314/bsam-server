package reboot

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func RebootPOST(c *gin.Context) {
	secret := c.Query("secret")

	if secret != os.Getenv("JWT_SECRET") {
		c.String(
			http.StatusUnauthorized,
			"Unauthorized",
		)
		return
	}

	c.String(
		http.StatusOK,
		"run os exit. will soon reboot",
	)

	os.Exit(0)
}
