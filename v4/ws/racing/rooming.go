package racing

import (
	"time"
)

//nolint:gochecknoglobals
var rooms = make(map[string]*Hub)

// AutoRooming creates a race room automatically.
func AutoRooming() {
	ticker := time.NewTicker(AutoRoomingInterval)
	defer ticker.Stop()

	for {
		scale()
		<-ticker.C
	}
}

// scale scales race room.
func scale() {
	assocs := getAssociationIDs()

	for _, id := range assocs {
		if _, ok := rooms[id]; !ok {
			rooms[id] = NewHub(id)
			go rooms[id].Run()
		}
	}
}

// getAssociations returns association IDs.
func getAssociationIDs() []string {
	return []string{"sailing-ise", "hogehoge", "piyopiyo"}
}
