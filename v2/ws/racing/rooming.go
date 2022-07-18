package racing

import (
	"time"
)

var (
	rooms = make(map[string]*Hub)
)

func AutoRooming() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	races := []string{"hogehoge", "piyopiyo"}

	for {
		scale(races)
		<-ticker.C
	}
}

func scale(races []string) {
	for _, id := range races {
		if _, exist := rooms[id]; !exist {
			rooms[id] = NewHub(id)
			go rooms[id].Run()
		}
	}
}
