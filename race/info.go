package race

import (
	"sailing-assist-mie-api/abort"
	"sailing-assist-mie-api/bsamdb"
	"sailing-assist-mie-api/inspector"
	"sailing-assist-mie-api/message"
	"time"

	"github.com/gin-gonic/gin"
)

type InfoPUTJSON struct {
	Name       string `json:"name"`
	StartAt    time.Time
	StartAtStr string `json:"start_at"`
	EndAt      time.Time
	EndAtStr   string   `json:"end_at"`
	PointA     string   `json:"point_a"`
	PointB     string   `json:"point_b"`
	PointC     string   `json:"point_c"`
	Athletes   []string `json:"athletes"`
	Memo       string   `json:"memo"`
	ImageUrl   string   `json:"image_url"`
	Holding    *bool    `json:"is_holding"`
}

// infoPUT is /race/:id PUT request handler.
func infoPUT(c *gin.Context) {
	ins := inspector.Inspector{Request: c.Request}
	raceId := c.Param("id")

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

	// Check convertable a timestamp value to time.Time.
	if json.StartAtStr != "" {
		json.StartAt, err = inspector.ParseTimestamp(json.StartAtStr)
		if err != nil {
			abort.BadRequest(c, message.NotMeetAllRequest)
			return
		}
	}

	if json.EndAtStr != "" {
		json.EndAt, err = inspector.ParseTimestamp(json.EndAtStr)
		if err != nil {
			abort.BadRequest(c, message.NotMeetAllRequest)
			return
		}
	}

	db, err := bsamdb.Open()
	if err != nil {
		panic(err)
	}
	defer db.DB.Close()

	// Check already stored this id.
	exist, err := db.IsExist("races", "id", raceId)
	if err != nil {
		panic(err)
	}

	// Update if already stored.
	if exist {
		err = update(&db, &json, raceId)
		if err != nil {
			switch err {
			case bsamdb.ErrRecordNotFound:
				abort.NotFound(c, message.RaceNotFound)
			default:
				panic(err)
			}
		}

	} else {
		abort.NotFound(c, message.RaceNotFound)
		return
	}
}

// Update updates to new data.
func update(db *bsamdb.DbInfo, json *InfoPUTJSON, raceId string) error {
	// Records
	data := []bsamdb.Field{}

	if json.Name != "" {
		data = append(data, bsamdb.Field{
			Column: "name",
			Value:  json.Name,
		})
	}

	if json.StartAtStr != "" {
		data = append(data, bsamdb.Field{
			Column: "start_at",
			Value:  json.StartAt,
		})
	}

	if json.EndAtStr != "" {
		data = append(data, bsamdb.Field{
			Column: "end_at",
			Value:  json.EndAt,
		})
	}

	if json.PointA != "" {
		data = append(data, bsamdb.Field{
			Column: "point_a",
			Value:  json.PointA,
		})
	}

	if json.PointB != "" {
		data = append(data, bsamdb.Field{
			Column: "point_b",
			Value:  json.PointB,
		})
	}

	if json.PointC != "" {
		data = append(data, bsamdb.Field{
			Column: "point_c",
			Value:  json.PointC,
		})
	}

	if len(json.Athletes) > 0 {
		data = append(data, bsamdb.Field{
			Column: "athlete",
			Value:  json.Athletes,
		})
	}

	if json.Memo != "" {
		data = append(data, bsamdb.Field{
			Column: "memo",
			Value:  json.Memo,
		})
	}

	if json.ImageUrl != "" {
		data = append(data, bsamdb.Field{
			Column: "image_url",
			Value:  json.ImageUrl,
		})
	}

	if json.Holding != nil {
		data = append(data, bsamdb.Field{
			Column: "is_holding",
			Value:  json.Holding,
		})
	}

	if len(data) > 0 {
		_, err := db.Update(
			"races",
			"id",
			raceId,
			data,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
