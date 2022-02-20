package device

import (
	"fmt"
	"sailing-assist-mie-api/abort"
	"sailing-assist-mie-api/bsamdb"
	"sailing-assist-mie-api/inspector"
	"sailing-assist-mie-api/message"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type InfoPOSTJSON struct {
	Name      string  `json:"name" binding:"required"`
	Model     int     `json:"model" binding:"required"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type InfoPUTJSON struct {
	Name      string  `json:"name"`
	Model     int     `json:"model"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// infoPOST is /device/:imei POST request handler.
func infoPOST(c *gin.Context) {
	ins := inspector.Inspector{Request: c.Request}
	imei := c.Param("imei")

	// Only JSON.
	if !ins.IsJSON() {
		abort.BadRequest(c, message.OnlyJSON)
		return
	}

	var json InfoPOSTJSON

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

	// Check already stored this imei.
	exist, err := db.IsExist("devices", "imei", imei)
	if err != nil {
		panic(err)
	}

	// Create if not stored.
	if !exist {
		err = create(&db, &json, imei)
		if err != nil {
			panic(err)
		}
	} else {
		abort.Conflict(c, message.AlreadyExisted)
		return
	}
}

// infoPUT is /device/:imei PUT request handler.
func infoPUT(c *gin.Context) {
	ins := inspector.Inspector{Request: c.Request}
	imei := c.Param("imei")

	// Only JSON.
	if !ins.IsJSON() {
		fmt.Println("1 out")
		abort.BadRequest(c, message.OnlyJSON)
		return
	}

	var json InfoPUTJSON

	// Check all of the require field is not blanked.
	err := c.ShouldBindBodyWith(&json, binding.JSON)
	if err != nil {
		fmt.Println("2 out")
		abort.BadRequest(c, message.NotMeetAllRequest)
		return
	}

	db, err := bsamdb.Open()
	if err != nil {
		panic(err)
	}
	defer db.DB.Close()

	// Check already stored this imei.
	exist, err := db.IsExist("devices", "imei", imei)
	if err != nil {
		panic(err)
	}

	// Update if already stored.
	if exist {
		err = update(&db, &json, imei)
		if err != nil {
			switch err {
			case bsamdb.ErrRecordNotFound:
				abort.NotFound(c, message.DeviceNotFound)
			default:
				panic(err)
			}
		}

	} else {
		var newJson InfoPOSTJSON

		// Check all of the require field is not blanked.
		err := c.ShouldBindBodyWith(&newJson, binding.JSON)
		if err != nil {
			fmt.Println("3 out")
			abort.BadRequest(c, message.NotMeetAllRequest)
			return
		}

		err = create(&db, &newJson, imei)
		if err != nil {
			panic(err)
		}
	}
}

// Create stores new device data.
func create(db *bsamdb.DbInfo, json *InfoPOSTJSON, imei string) error {
	// Records
	data := []bsamdb.Field{
		{Column: "imei", Value: imei},
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

	_, err := db.Insert(
		"devices",
		data,
	)
	if err != nil {
		return err
	}

	return nil
}

// Update updates to new data.
func update(db *bsamdb.DbInfo, json *InfoPUTJSON, imei string) error {
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

	if len(data) > 0 {
		_, err := db.Update(
			"devices",
			"imei",
			imei,
			data,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
