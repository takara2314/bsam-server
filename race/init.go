package race

import (
	"sailing-assist-mie-api/bsamdb"
	"sailing-assist-mie-api/utils"
)

func init() {
	// Connect to the database and insert such data.
	db, err := bsamdb.Open()
	if err != nil {
		panic(err)
	}
	defer db.DB.Close()

	// Update the database to reset registered athletes info.
	_, err = db.UpdateAll(
		"races",
		[]bsamdb.Field{{
			Column: "athlete",
			Value2d: utils.StrSliceToAnySlice(
				[]string{},
			),
		}},
	)
}
