package device

import (
	"sailing-assist-mie-api/abort"
	"sailing-assist-mie-api/bsamdb"
	"sailing-assist-mie-api/inspector"
	"sailing-assist-mie-api/message"

	"github.com/gin-gonic/gin"
)

type DevicesPOSTJSON struct {
	IMEI      string  `json:"imei" binding:"required"`
	Name      string  `json:"name" binding:"required"`
	Model     int     `json:"model" binding:"required"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type DevicesPUTJSON struct {
	IMEI      string  `json:"imei" binding:"required"`
	Name      string  `json:"name"`
	Model     int     `json:"model"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// DevicesPOST is /devices POST request handler.
func DevicesPOST(c *gin.Context) {
	ins := inspector.Inspector{Request: c.Request}

	// Only JSON.
	if !ins.IsJSON() {
		abort.BadRequest(c, message.OnlyJSON)
		return
	}

	var json DevicesPOSTJSON

	// Check all of the require field is not blanked.
	err := c.ShouldBindJSON(&json)
	if err != nil {
		abort.BadRequest(c, message.NotMeetAllRequest)
		return
	}

	// Records
	data := []bsamdb.Field{
		{Column: "imei", Value: json.IMEI},
		{Column: "name", Value: json.Name},
		{Column: "model", Value: json.Model},
	}

	if json.Latitude != 0.0 {
		data = append(data, bsamdb.Field{
			Column: "latitude",
			Value:  json.Latitude,
		})
	}

	if json.Longitude != 0.0 {
		data = append(data, bsamdb.Field{
			Column: "longitude",
			Value:  json.Longitude,
		})
	}

	// Connect to the database and insert such data.
	db, err := bsamdb.Open()
	if err != nil {
		panic(err)
	}
	defer db.DB.Close()

	_, err = db.Insert(
		"devices",
		data,
	)
	if err != nil {
		panic(err)
	}
}

// DevicesPUT is /devices PUT request handler.
func DevicesPUT(c *gin.Context) {
	ins := inspector.Inspector{Request: c.Request}

	// Only JSON.
	if !ins.IsJSON() {
		abort.BadRequest(c, message.OnlyJSON)
		return
	}

	var json DevicesPUTJSON

	// Check all of the require field is not blanked.
	err := c.ShouldBindJSON(&json)
	if err != nil {
		abort.BadRequest(c, message.NotMeetAllRequest)
		return
	}

	// Records
	data := []bsamdb.Field{}

	if json.Name != "" {
		data = append(data, bsamdb.Field{
			Column: "name",
			Value:  json.Name,
		})
	}

	if json.Model != 0 {
		data = append(data, bsamdb.Field{
			Column: "model",
			Value:  json.Model,
		})
	}

	if json.Latitude != 0.0 {
		data = append(data, bsamdb.Field{
			Column: "latitude",
			Value:  json.Latitude,
		})
	}

	if json.Longitude != 0.0 {
		data = append(data, bsamdb.Field{
			Column: "longitude",
			Value:  json.Longitude,
		})
	}

	if len(data) != 0 {
		// Connect to the database and update such data.
		db, err := bsamdb.Open()
		if err != nil {
			panic(err)
		}
		defer db.DB.Close()

		_, err = db.Update(
			"devices",
			"imei",
			json.IMEI,
			data,
		)

		if err != nil {
			switch err {
			case bsamdb.ErrRecordNotFound:
				abort.NotFound(c, message.DeviceNotFound)
			default:
				panic(err)
			}
		}
	}
}
