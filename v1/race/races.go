package race

import (
	"bsam-server/utils"
	"bsam-server/v1/abort"
	"bsam-server/v1/bsamdb"
	"bsam-server/v1/inspector"
	"bsam-server/v1/message"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type RaceInfo struct {
	ID       *string    `json:"id"`
	Name     *string    `json:"name"`
	StartAt  *time.Time `json:"start_at"`
	EndAt    *time.Time `json:"end_at"`
	PointA   *string    `json:"point_a"`
	PointB   *string    `json:"point_b"`
	PointC   *string    `json:"point_c"`
	Athletes []string   `json:"athletes"`
	Memo     *string    `json:"memo"`
	ImageURL *string    `json:"image_url"`
	Holding  *bool      `json:"is_holding"`
}

type RaceGETJSON struct {
	Status string     `json:"status"`
	Races  []RaceInfo `json:"races"`
}

type RacePOSTJSON struct {
	Name       string `json:"name" binding:"required"`
	StartAt    time.Time
	StartAtStr string `json:"start_at" binding:"required"`
	EndAt      time.Time
	EndAtStr   string   `json:"end_at" binding:"required"`
	PointA     string   `json:"point_a"`
	PointB     string   `json:"point_b"`
	PointC     string   `json:"point_c"`
	Athletes   []string `json:"athletes"`
	Memo       string   `json:"memo"`
	ImageURL   string   `json:"image_url"`
	Holding    *bool    `json:"is_holding" binding:"required"`
}

// RacesGET is /races GET request handler.
func RacesGET(c *gin.Context) {
	// ins := inspector.Inspector{Request: c.Request}

	// Connect to the database and insert such data.
	db, err := bsamdb.Open()
	if err != nil {
		panic(err)
	}
	defer db.DB.Close()

	races, err := fetchAll(&db, "")
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, RaceGETJSON{
		Status: "OK",
		Races:  races,
	})
}

// RacesPOST is /races POST request handler.
func RacesPOST(c *gin.Context) {
	ins := inspector.Inspector{Request: c.Request}

	// Only JSON.
	if !ins.IsJSON() {
		abort.BadRequest(c, message.OnlyJSON)
		return
	}

	var json RacePOSTJSON

	// Check all of the require field is not blanked.
	err := c.ShouldBindJSON(&json)
	if err != nil {
		abort.BadRequest(c, message.NotMeetAllRequest)
		return
	}

	// Check convertable a timestamp value to time.Time.
	json.StartAt, err = inspector.ParseTimestamp(json.StartAtStr)
	if err != nil {
		abort.BadRequest(c, message.NotMeetAllRequest)
		return
	}

	json.EndAt, err = inspector.ParseTimestamp(json.EndAtStr)
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

	err = create(&db, &json)
	if err != nil {
		panic(err)
	}
}

// fetchAll fetches all of rows in this group.
func fetchAll(db *bsamdb.DbInfo, groupID string) ([]RaceInfo, error) {
	races := make([]RaceInfo, 0)
	data := make([]bsamdb.Field, 0)

	if groupID != "" {
		data = append(
			data,
			bsamdb.Field{Column: "group_id", Value: groupID},
		)
	}

	rows, err := db.Select(
		"races",
		data,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		info := RaceInfo{}
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
		startAtLocal := (*info.StartAt).Local()
		info.StartAt = &startAtLocal
		endAtLocal := (*info.EndAt).Local()
		info.EndAt = &endAtLocal

		if err != nil {
			return nil, err
		}

		races = append(races, info)
	}

	return races, nil
}

// Create stores new device data.
func create(db *bsamdb.DbInfo, json *RacePOSTJSON) error {
	// Records
	data := []bsamdb.Field{
		{Column: "name", Value: json.Name},
		{Column: "start_at", Value: json.StartAt},
		{Column: "end_at", Value: json.EndAt},
		{Column: "is_holding", Value: json.Holding},
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
			Column:  "athlete",
			Value2d: utils.StrSliceToAnySlice(json.Athletes),
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

	_, err := db.Insert(
		"races",
		data,
	)
	if err != nil {
		return err
	}

	return nil
}
