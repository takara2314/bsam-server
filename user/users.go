package user

import (
	"net/http"
	"sailing-assist-mie-api/abort"
	"sailing-assist-mie-api/inspector"
	"sailing-assist-mie-api/message"

	"github.com/gin-gonic/gin"
)

type UsersPOSTJSON struct {
	LoginName   string  `json:"login_name" binding:"required"`
	DisplayName string  `json:"display_name" binding:"required"`
	Password    string  `json:"password" binding:"required"`
	GroupId     string  `json:"group_id"`
	Role        string  `json:"role"`
	DeviceIMEI  string  `json:"device"`
	SailNum     int     `json:"sail_num"`
	CourseLimit float32 `json:"course_limit"`
	Image       string  `json:"image"`
	Note        string  `json:"note"`
}

// UsersPOST is /users POST request handler.
func UsersPOST(c *gin.Context) {
	ins := inspector.Inspector{Request: c.Request}

	if mes := ins.HasToken(); mes != "" {
		abort.Unauthorized(c, mes)
		return
	}

	// Require "admin.user.users.create" or "developer.user.users.create".
	if !ins.HasPermission([]string{
		"admin.user.users.create",
		"developer.user.users.create",
	}) {
		abort.Forbidden(c, message.CannotUseAPI)
		return
	}

	// Only JSON
	if !ins.IsJSON() {
		abort.BadRequest(c, message.OnlyJSON)
		return
	}

	var json UsersPOSTJSON

	err := c.ShouldBindJSON(&json)
	if err != nil {
		abort.BadRequest(c, message.NotMeetAllRequest)
		return
	}

	// Not developer cannot add outside user. Only inside user.
	if json.GroupId != "" {
		if !ins.IsSameGroup(json.GroupId) && !ins.HasPermission([]string{"developer.user.users.create"}) {
			abort.Forbidden(c, message.CannotAddOutSideUser)
			return
		}
	}

	c.String(http.StatusOK, "Hello "+json.DisplayName)
}
