package race

import (
	"net/http"
	"sailing-assist-mie-api/abort"
	"sailing-assist-mie-api/bsamdb"
	"sailing-assist-mie-api/inspector"
	"sailing-assist-mie-api/message"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type InfoGETJSON struct {
	Status string   `json:"status"`
	Race   RaceInfo `json:"race"`
}

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
	ImageURL   string   `json:"image_url"`
	Holding    *bool    `json:"is_holding"`
}

// infoGET is /race/:id GET request handler.
func infoGET(c *gin.Context) {
	// ins := inspector.Inspector{Request: c.Request}
	raceID := c.Param("id")

	// Connect to the database and insert such data.
	db, err := bsamdb.Open()
	if err != nil {
		panic(err)
	}
	defer db.DB.Close()

	// Check already stored this id.
	exist, err := db.IsExist("races", "id", raceID)
	if err != nil {
		panic(err)
	}

	// Update if already stored.
	if exist {
		race, err := fetch(&db, raceID)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, InfoGETJSON{
			Status: "OK",
			Race:   race,
		})
	} else {
		abort.NotFound(c, message.RaceNotFound)
		return
	}
}

// infoPUT is /race/:id PUT request handler.
func infoPUT(c *gin.Context) {
	ins := inspector.Inspector{Request: c.Request}
	raceID := c.Param("id")

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
	exist, err := db.IsExist("races", "id", raceID)
	if err != nil {
		panic(err)
	}

	// Update if already stored.
	if exist {
		err = update(&db, &json, raceID)
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

// fetch fetches rows in this group.
func fetch(db *bsamdb.DbInfo, raceID string) (RaceInfo, error) {
	rows, err := db.Select(
		"races",
		[]bsamdb.Field{
			{Column: "id", Value: raceID},
		},
	)
	if err != nil {
		return RaceInfo{}, err
	}
	defer rows.Close()

	info := RaceInfo{}
	rows.Next()
	err = rows.Scan(
		&info.ID,
		&info.Name,
		&info.StartAt,
		&info.EndAt,
		&info.PointA,
		&info.PointB,
		&info.PointC,
		pq.Array(&info.Athletes),
		&info.Memo,
		&info.ImageURL,
		&info.Holding,
	)
	if err != nil {
		return RaceInfo{}, err
	}

	return info, nil
}

// Update updates to new data.
func update(db *bsamdb.DbInfo, json *InfoPUTJSON, raceID string) error {
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

	if json.ImageURL != "" {
		data = append(data, bsamdb.Field{
			Column: "image_url",
			Value:  json.ImageURL,
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
			raceID,
			data,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
