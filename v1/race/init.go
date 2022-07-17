package race

import (
	"bsam-server/utils"
	"bsam-server/v1/bsamdb"
	"fmt"
)

func init() {
	// Connect to the database and insert such data.
	db, err := bsamdb.Open()
	if err != nil {
		panic(err)
	}
	defer db.DB.Close()

	fmt.Println("今から実行するで！")

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

	if err != nil {
		fmt.Println(err)
	}
}
