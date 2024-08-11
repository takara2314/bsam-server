package geolocationhub

import (
	"time"

	"cloud.google.com/go/firestore"
)

type Hub struct {
	AssocID         string
	FirestoreClient *firestore.Client
}

type Geolocation struct {
	DeviceID              string
	Latitude              float64
	Longitude             float64
	AltitudeMeter         float64
	AccuracyMeter         float64
	AltitudeAccuracyMeter float64
	Heading               float64
	SpeedMeterPerSec      float64
	RecordedAt            time.Time
}

func NewHub(
	AssocID string,
	firestoreClient *firestore.Client,
) *Hub {
	return &Hub{
		AssocID:         AssocID,
		FirestoreClient: firestoreClient,
	}
}
