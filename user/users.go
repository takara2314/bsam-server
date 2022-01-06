package user

import (
	"net/http"
	"sailing-assist-mie-api/abort"
	"sailing-assist-mie-api/inspector"

	"github.com/gin-gonic/gin"
)

type UsersPOSTJSON struct {
	LoginName   string  `json:"login_name" binding:"required"`
	DisplayName string  `json:"display_name" binding:"required"`
	Password    string  `json:"password" binding:"required"`
	UserType    string  `json:"user_type"`
	DeviceIMEI  string  `json:"device"`
	SailNum     int     `json:"sail_num"`
	CourseLimit float32 `json:"course_limit"`
	Image       string  `json:"image"`
	Note        string  `json:"note"`
}

// UsersPOST is /users POST request handler
func UsersPOST(c *gin.Context) {
	ins := inspector.Inspector{Context: c}

	if !ins.HasPermission("admin.user.users.create") {
		abort.Forbidden(c, "You cannot use this API.")
		return
	}

	var json UsersPOSTJSON

	err := c.ShouldBindJSON(&json)
	if err != nil {
		abort.BadRequest(c, "This request does not meet all of the required elements.")
		return
	}

	c.String(http.StatusOK, "Hello "+json.DisplayName)
}
