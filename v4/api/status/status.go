package status

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

func GET(c *gin.Context) {
	var m runtime.MemStats

	runtime.ReadMemStats(&m)

	alloc := float64(m.Alloc) / (1024 * 1024)

	c.String(
		http.StatusOK,
		fmt.Sprintf("Using: %f MB", alloc),
	)
}
