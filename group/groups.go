package group

import (
	"sailing-assist-mie-api/abort"
	"sailing-assist-mie-api/bsamdb"
	"sailing-assist-mie-api/inspector"
	"sailing-assist-mie-api/message"

	"github.com/gin-gonic/gin"
)

type GroupPOSTJSON struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// GroupsPOST is /groups POST request handler.
func GroupsPOST(c *gin.Context) {
	ins := inspector.Inspector{Request: c.Request}

	// Only JSON.
	if !ins.IsJSON() {
		abort.BadRequest(c, message.OnlyJSON)
		return
	}

	var json GroupPOSTJSON

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

	err = create(&db, &json)
	if err != nil {
		panic(err)
	}
}

// Create stores new device data.
func create(db *bsamdb.DbInfo, json *GroupPOSTJSON) error {
	// Records
	data := []bsamdb.Field{
		{Column: "name", Value: json.Name},
	}

	if json.Description != "" {
		data = append(data, bsamdb.Field{
			Column: "description",
			Value:  json.Description,
		})
	}

	_, err := db.Insert(
		"groups",
		data,
	)
	if err != nil {
		return err
	}

	return nil
}
