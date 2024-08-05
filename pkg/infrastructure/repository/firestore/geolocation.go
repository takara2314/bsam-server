package firestore

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/samber/oops"
)

type Geolocation struct {
	ID                    string    `firestore:"-"`
	Latitude              float64   `firestore:"latitude"`
	Longitude             float64   `firestore:"longitude"`
	AltitudeMeter         float64   `firestore:"altitudeMeter"`
	AccuracyMeter         float64   `firestore:"accuracyMeter"`
	AltitudeAccuracyMeter float64   `firestore:"altitudeAccuracyMeter"`
	Heading               float64   `firestore:"heading"`
	SpeedMeterPerSec      float64   `firestore:"speedMeterPerSec"`
	UpdatedAt             time.Time `firestore:"updatedAt"`
}

func SetGeolocation(
	ctx context.Context,
	client *firestore.Client,
	id string,
	lat float64,
	lng float64,
	altMeter float64,
	accMeter float64,
	altAccMeter float64,
	heading float64,
	speedMeterPerSec float64,
	updatedAt time.Time,
) error {
	_, err := client.Collection("geolocations").Doc(id).Set(ctx, Geolocation{
		ID:                    id,
		Latitude:              lat,
		Longitude:             lng,
		AltitudeMeter:         altMeter,
		AccuracyMeter:         accMeter,
		AltitudeAccuracyMeter: altAccMeter,
		Heading:               heading,
		SpeedMeterPerSec:      speedMeterPerSec,
		UpdatedAt:             updatedAt,
	})
	if err != nil {
		return oops.
			In("repository.AddAssoc").
			Wrapf(err, "failed to add assoc")
	}

	return nil
}

func FetchGeolocationByID(ctx context.Context, client *firestore.Client, id string) (*Geolocation, error) {
	doc, err := client.Collection("geolocations").Doc(id).Get(ctx)
	if err != nil {
		return nil, oops.
			In("repository.FetchGeolocationByID").
			Wrapf(err, "failed to fetch geolocation")
	}

	var loc Geolocation
	err = doc.DataTo(&loc)
	if err != nil {
		return nil, oops.
			In("repository.FetchGeolocationByID").
			Wrapf(err, "failed to convert data to user")
	}

	loc.ID = doc.Ref.ID

	return &loc, err
}
