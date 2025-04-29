package status

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

const OneKB = 1024
const OneMB = OneKB * OneKB

func GET(c *gin.Context) {
	var m runtime.MemStats

	runtime.ReadMemStats(&m)

	alloc := float64(m.Alloc) / OneMB

	c.String(
		http.StatusOK,
		fmt.Sprintf("Using: %f MB", alloc),
	)
}
