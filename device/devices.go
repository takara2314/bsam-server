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

// DevicesPOST is /devices POST request handler.
func DevicesPOST(c *gin.Context) {
	ins := inspector.Inspector{Request: c.Request}

	// OnlyJSON
	if !ins.IsJSON() {
		abort.BadRequest(c, message.OnlyJSON)
		return
	}

	var json DevicesPOSTJSON

	err := c.ShouldBindJSON(&json)
	if err != nil {
		abort.BadRequest(c, message.NotMeetAllRequest)
		return
	}

	data := []bsamdb.Field{
		{Column: "imei", Value: json.IMEI},
		{Column: "name", Value: json.Name},
		{Column: "model", Value: json.Model},
	}

	if json.Latitude != 0 {
		data = append(data, bsamdb.Field{
			Column: "latitude",
			Value:  json.Latitude,
		})
	}

	if json.Longitude != 0 {
		data = append(data, bsamdb.Field{
			Column: "longitude",
			Value:  json.Longitude,
		})
	}

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
