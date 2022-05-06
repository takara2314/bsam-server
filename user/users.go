package user

import (
	"net/http"
	"sailing-assist-mie-api/abort"
	"sailing-assist-mie-api/bsamdb"
	"sailing-assist-mie-api/inspector"
	"sailing-assist-mie-api/message"

	"github.com/gin-gonic/gin"
)

type UserGETJSON struct {
	Status string     `json:"status"`
	Users  []UserInfo `json:"users"`
}

type UserPOSTJSON struct {
	LoginID     string  `json:"login_id" binding:"required"`
	DisplayName string  `json:"display_name" binding:"required"`
	Password    string  `json:"password" binding:"required"`
	GroupID     string  `json:"group_id" binding:"required"`
	Role        string  `json:"role" binding:"required"`
	DeviceID    string  `json:"device_id"`
	SailNum     int     `json:"sail_num"`
	CourseLimit float32 `json:"course_limit"`
	ImageURL    string  `json:"image_url"`
	Note        string  `json:"note"`
}

// UsersGET is /users GET request handler.
func UsersGET(c *gin.Context) {
	// ins := inspector.Inspector{Request: c.Request}

	// Connect to the database and insert such data.
	db, err := bsamdb.Open()
	if err != nil {
		panic(err)
	}
	defer db.DB.Close()

	users, err := fetchAll(&db, c.Query("role"))
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, UserGETJSON{
		Status: "OK",
		Users:  users,
	})
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
	exist, err := db.IsExist("users", "login_id", json.LoginID)
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

// fetchAll fetches all of rows.
//   selector: role
func fetchAll(db *bsamdb.DbInfo, role string) ([]UserInfo, error) {
	users := make([]UserInfo, 0)
	data := make([]bsamdb.Field, 0)

	if role != "" {
		data = append(
			data,
			bsamdb.Field{Column: "role", Value: role},
		)
	}

	rows, err := db.Select(
		"users",
		data,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var password string

	for rows.Next() {
		info := UserInfo{}
		err = rows.Scan(
			&info.UserID,
			&info.LoginID,
			&info.DisplayName,
			&password,
			&info.GroupID,
			&info.Role,
			&info.DeviceID,
			&info.SailNum,
			&info.CourseLimit,
			&info.ImageURL,
			&info.Note,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, info)
	}

	return users, nil
}

// Create stores new device data.
func create(db *bsamdb.DbInfo, json *UserPOSTJSON) error {
	// Records
	data := []bsamdb.Field{
		{Column: "login_id", Value: json.LoginID},
		{Column: "display_name", Value: json.DisplayName},
		{Column: "password", Value: json.Password, ToHash: true},
		{Column: "group_id", Value: json.GroupID},
		{Column: "role", Value: json.Role},
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

	_, err := db.Insert(
		"users",
		data,
	)
	if err != nil {
		return err
	}

	return nil
}
