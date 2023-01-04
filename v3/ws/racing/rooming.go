package racing

import (
	"time"
)

var (
	rooms = make(map[string]*Hub)
)

// AutoRooming creates a race room automatically.
func AutoRooming() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	races := []string{"hogehoge", "piyopiyo", "3ae8c214-eb72-481c-b110-8e8f32ecf02d"}

	for {
		scale(races)
		<-ticker.C
	}
}

// scale scales race room.
func scale(races []string) {
	for _, id := range races {
		if _, exist := rooms[id]; !exist {
			rooms[id] = NewHub(id)
			go rooms[id].Run()
		}
	}
}
