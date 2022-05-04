package user

import (
	"database/sql"
	"net/http"

	"sailing-assist-mie-api/abort"
	"sailing-assist-mie-api/bsamdb"
	"sailing-assist-mie-api/inspector"
	"sailing-assist-mie-api/message"

	"github.com/gin-gonic/gin"
)

type UserInfoNullable struct {
	UserID      string
	LoginID     string
	DisplayName string
	Password    string
	GroupID     string
	Role        string
	DeviceID    sql.NullString
	SailNum     sql.NullInt16
	CourseLimit sql.NullFloat64
	ImageURL    sql.NullString
	Note        sql.NullString
}

type UserInfo struct {
	UserID      string  `json:"user_id"`
	LoginID     string  `json:"login_id"`
	DisplayName string  `json:"display_name"`
	GroupID     string  `json:"group_id"`
	Role        string  `json:"role"`
	DeviceID    string  `json:"device_id"`
	SailNum     int     `json:"sail_num"`
	CourseLimit float32 `json:"course_limit"`
	ImageURL    string  `json:"image_url"`
	Note        string  `json:"note"`
}

type InfoGETResponse struct {
	Status string   `json:"status"`
	Info   UserInfo `json:"info"`
}

type InfoPUTJSON struct {
	LoginID     string  `json:"login_id"`
	DisplayName string  `json:"display_name"`
	Password    string  `json:"password"`
	GroupID     string  `json:"group_id"`
	Role        string  `json:"role"`
	DeviceID    string  `json:"device_id"`
	SailNum     int     `json:"sail_num"`
	CourseLimit float32 `json:"course_limit"`
	ImageURL    string  `json:"image_url"`
	Note        string  `json:"note"`
}

// infoGET is /user/:id GET request handler.
func infoGET(c *gin.Context) {
	// ins := inspector.Inspector{Request: c.Request}
	userID := c.Param("id")

	db, err := bsamdb.Open()
	if err != nil {
		panic(err)
	}
	defer db.DB.Close()

	// Check already stored this id.
	exist, err := db.IsExist("users", "id", userID)
	if err != nil {
		panic(err)
	}

	if !exist {
		abort.NotFound(c, message.UserNotFound)
		return
	}

	rows, err := db.Select(
		"users",
		[]bsamdb.Field{
			{Column: "id", Value: userID},
		},
	)
	if err != nil {
		panic(err)
	}

	tmp := UserInfoNullable{}

	rows.Next()
	err = rows.Scan(
		&tmp.UserID,
		&tmp.LoginID,
		&tmp.DisplayName,
		&tmp.Password,
		&tmp.GroupID,
		&tmp.Role,
		&tmp.DeviceID,
		&tmp.SailNum,
		&tmp.CourseLimit,
		&tmp.ImageURL,
		&tmp.Note,
	)
	if err != nil {
		panic(err)
	}

	info := UserInfo{
		UserID:      tmp.UserID,
		LoginID:     tmp.LoginID,
		DisplayName: tmp.DisplayName,
		GroupID:     tmp.GroupID,
		Role:        tmp.Role,
		DeviceID:    tmp.DeviceID.String,
		SailNum:     int(tmp.SailNum.Int16),
		CourseLimit: float32(tmp.CourseLimit.Float64),
		ImageURL:    tmp.ImageURL.String,
		Note:        tmp.Note.String,
	}

	c.JSON(http.StatusOK, InfoGETResponse{
		Status: "OK",
		Info:   info,
	})
}

// infoPUT is /user/:id PUT request handler.
func infoPUT(c *gin.Context) {
	ins := inspector.Inspector{Request: c.Request}
	userID := c.Param("id")

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
	exist, err := db.IsExist("users", "id", userID)
	if err != nil {
		panic(err)
	}

	// Update if already stored.
	if exist {
		if json.LoginID != "" {
			// Check already stored this login_id.
			exist, err := db.IsExistNotIt("users", "id", userID, "login_id", json.LoginID)
			if err != nil {
				panic(err)
			}

			if exist {
				abort.Conflict(c, message.AlreadyExisted)
				return
			}
		}

		err = update(&db, &json, userID)
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
func update(db *bsamdb.DbInfo, json *InfoPUTJSON, userID string) error {
	// Records
	data := []bsamdb.Field{}

	if json.LoginID != "" {
		data = append(data, bsamdb.Field{
			Column: "login_id",
			Value:  json.LoginID,
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

	if json.GroupID != "" {
		data = append(data, bsamdb.Field{
			Column: "group_id",
			Value:  json.GroupID,
		})
	}

	if json.Role != "" {
		data = append(data, bsamdb.Field{
			Column: "role",
			Value:  json.Role,
		})
	}

	if json.DeviceID != "" {
		data = append(data, bsamdb.Field{
			Column: "device_id",
			Value:  json.DeviceID,
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

	if json.ImageURL != "" {
		data = append(data, bsamdb.Field{
			Column: "image_url",
			Value:  json.ImageURL,
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
			userID,
			data,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
