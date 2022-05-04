package group

import (
	"sailing-assist-mie-api/abort"
	"sailing-assist-mie-api/bsamdb"
	"sailing-assist-mie-api/inspector"
	"sailing-assist-mie-api/message"

	"github.com/gin-gonic/gin"
)

type InfoPUTJSON struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// infoPUT is /group/:id PUT request handler.
func infoPUT(c *gin.Context) {
	ins := inspector.Inspector{Request: c.Request}
	groupID := c.Param("id")

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
	exist, err := db.IsExist("groups", "id", groupID)
	if err != nil {
		panic(err)
	}

	// Update if already stored.
	if exist {
		err = update(&db, &json, groupID)
		if err != nil {
			switch err {
			case bsamdb.ErrRecordNotFound:
				abort.NotFound(c, message.GroupNotFound)
			default:
				panic(err)
			}
		}

	} else {
		abort.NotFound(c, message.GroupNotFound)
		return
	}
}

// Update updates to new data.
func update(db *bsamdb.DbInfo, json *InfoPUTJSON, groupID string) error {
	// Records
	data := []bsamdb.Field{}

	if json.Name != "" {
		data = append(data, bsamdb.Field{
			Column: "name",
			Value:  json.Name,
		})
	}

	if json.Description != "" {
		data = append(data, bsamdb.Field{
			Column: "description",
			Value:  json.Description,
		})
	}

	if len(data) > 0 {
		_, err := db.Update(
			"groups",
			"id",
			groupID,
			data,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
