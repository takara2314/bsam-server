package user

import (
	"sailing-assist-mie-api/abort"
	"sailing-assist-mie-api/bsamdb"
	"sailing-assist-mie-api/inspector"
	"sailing-assist-mie-api/message"

	"github.com/gin-gonic/gin"
)

type InfoPUTJSON struct {
	LoginId     string  `json:"login_id"`
	DisplayName string  `json:"display_name"`
	Password    string  `json:"password"`
	GroupId     string  `json:"group_id"`
	Role        string  `json:"role"`
	DeviceId    string  `json:"device_id"`
	SailNum     int     `json:"sail_num"`
	CourseLimit float32 `json:"course_limit"`
	ImageUrl    string  `json:"image_url"`
	Note        string  `json:"note"`
}

// infoPUT is /user/:id PUT request handler.
func infoPUT(c *gin.Context) {
	ins := inspector.Inspector{Request: c.Request}
	userId := c.Param("id")

	// Only JSON.
	if !ins.IsJSON() {
		abort.BadRequest(c, message.OnlyJSON)
		return
	}

	var json InfoPUTJSON

	// Check all of the require field is not blanked.
	err := c.ShouldBindJSON(&json)
	if err != nil {
		abort.BadRequest(c, message.NotMeetAllRequest)
		return
	}

	db, err := bsamdb.Open()
	if err != nil {
		panic(err)
	}
	defer db.DB.Close()

	// Check already stored this id.
	exist, err := db.IsExist("users", "id", userId)
	if err != nil {
		panic(err)
	}

	// Update if already stored.
	if exist {
		if json.LoginId != "" {
			// Check already stored this login_id.
			exist, err := db.IsExistNotIt("users", "id", userId, "login_id", json.LoginId)
			if err != nil {
				panic(err)
			}

			if exist {
				abort.Conflict(c, message.AlreadyExisted)
				return
			}
		}

		err = update(&db, &json, userId)
		if err != nil {
			switch err {
			case bsamdb.ErrRecordNotFound:
				abort.NotFound(c, message.UserNotFound)
			default:
				panic(err)
			}
		}

	} else {
		abort.NotFound(c, message.UserNotFound)
		return
	}
}

// Update updates to new data.
func update(db *bsamdb.DbInfo, json *InfoPUTJSON, userId string) error {
	// Records
	data := []bsamdb.Field{}

	if json.LoginId != "" {
		data = append(data, bsamdb.Field{
			Column: "login_id",
			Value:  json.LoginId,
		})
	}

	if json.DisplayName != "" {
		data = append(data, bsamdb.Field{
			Column: "display_name",
			Value:  json.DisplayName,
		})
	}

	if json.Password != "" {
		data = append(data, bsamdb.Field{
			Column: "password",
			Value:  json.Password,
			ToHash: true,
		})
	}

	if json.GroupId != "" {
		data = append(data, bsamdb.Field{
			Column: "group_id",
			Value:  json.GroupId,
		})
	}

	if json.Role != "" {
		data = append(data, bsamdb.Field{
			Column: "role",
			Value:  json.Role,
		})
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

	if len(data) > 0 {
		_, err := db.Update(
			"users",
			"id",
			userId,
			data,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
