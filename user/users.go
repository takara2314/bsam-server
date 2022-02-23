package user

import (
	"sailing-assist-mie-api/abort"
	"sailing-assist-mie-api/bsamdb"
	"sailing-assist-mie-api/inspector"
	"sailing-assist-mie-api/message"

	"github.com/gin-gonic/gin"
)

type UserPOSTJSON struct {
	LoginId     string  `json:"login_id" binding:"required"`
	DisplayName string  `json:"display_name" binding:"required"`
	Password    string  `json:"password" binding:"required"`
	GroupId     string  `json:"group_id" binding:"required"`
	Role        string  `json:"role" binding:"required"`
	DeviceId    string  `json:"device_id"`
	SailNum     int     `json:"sail_num"`
	CourseLimit float32 `json:"course_limit"`
	ImageUrl    string  `json:"image_url"`
	Note        string  `json:"note"`
}

// UsersPOST is /users POST request handler.
func UsersPOST(c *gin.Context) {
	ins := inspector.Inspector{Request: c.Request}

	// Only JSON.
	if !ins.IsJSON() {
		abort.BadRequest(c, message.OnlyJSON)
		return
	}

	var json UserPOSTJSON

	// Check all of the require field is not blanked.
	err := c.ShouldBindJSON(&json)
	if err != nil {
		abort.BadRequest(c, message.NotMeetAllRequest)
		return
	}

	// Connect to the database and insert such data.
	db, err := bsamdb.Open()
	if err != nil {
		panic(err)
	}
	defer db.DB.Close()

	// Check already stored this login_id.
	exist, err := db.IsExist("users", "login_id", json.LoginId)
	if err != nil {
		panic(err)
	}

	if !exist {
		err = create(&db, &json)
		if err != nil {
			panic(err)
		}
	} else {
		abort.Conflict(c, message.AlreadyExisted)
		return
	}
}

// Create stores new device data.
func create(db *bsamdb.DbInfo, json *UserPOSTJSON) error {
	// Records
	data := []bsamdb.Field{
		{Column: "login_id", Value: json.LoginId},
		{Column: "display_name", Value: json.DisplayName},
		{Column: "password", Value: json.Password, ToHash: true},
		{Column: "group_id", Value: json.GroupId},
		{Column: "role", Value: json.Role},
	}

	if json.DeviceId != "" {
		data = append(data, bsamdb.Field{
			Column: "device_id",
			Value:  json.DeviceId,
		})
	}

	if json.SailNum != 0 {
		data = append(data, bsamdb.Field{
			Column: "sail_num",
			Value:  json.SailNum,
		})
	}

	if json.CourseLimit != 0.0 {
		data = append(data, bsamdb.Field{
			Column: "course_limit",
			Value:  json.CourseLimit,
		})
	}

	if json.ImageUrl != "" {
		data = append(data, bsamdb.Field{
			Column: "image_url",
			Value:  json.ImageUrl,
		})
	}

	if json.Note != "" {
		data = append(data, bsamdb.Field{
			Column: "note",
			Value:  json.Note,
		})
	}

	_, err := db.Insert(
		"users",
		data,
	)
	if err != nil {
		return err
	}

	return nil
}
