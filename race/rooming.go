package race

import (
	"fmt"
	"sailing-assist-mie-api/bsamdb"
	"time"
)

var (
	rooms = make(map[string]*Hub)
)

// AutoRooming auto-scaling to run racing instance
// from reserved race info once every 30s.
func AutoRooming() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	db, err := bsamdb.Open()
	if err != nil {
		panic(err)
	}
	defer db.DB.Close()

	for {
		scale(&db)
		<-ticker.C
	}
}

func scale(db *bsamdb.DbInfo) {
	races, err := fetchAll(db, "")
	if err != nil {
		panic(err)
	}

	// If not exist race instance, create a race instance and run.
	for _, race := range races {
		if _, exist := rooms[*race.Id]; !exist {
			fmt.Println(*race.Id, "ルームを作成します")
			rooms[*race.Id] = NewHub(*race.Id)
			go rooms[*race.Id].Run()
		}
	}
}
