package geolocationhub

import (
	"context"
	"log/slog"
	"time"

	"github.com/samber/oops"
	repoFirestore "github.com/takara2314/bsam-server/pkg/infrastructure/repository/firestore"
)

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

func (h *Hub) StoreGeolocation(
	ctx context.Context,
	deviceID string,
	lat float64,
	lng float64,
	altMeter float64,
	accMeter float64,
	altAccMeter float64,
	heading float64,
	speedMeterPerSec float64,
	recordedAt time.Time,
) error {
	geolocationID := h.AssocID + "_" + deviceID

	if err := repoFirestore.SetGeolocation(
		ctx,
		h.FirestoreClient,
		geolocationID,
		lat,
		lng,
		altMeter,
		accMeter,
		altAccMeter,
		heading,
		speedMeterPerSec,
		recordedAt,
	); err != nil {
		return oops.
			In("geolocationhub.StoreGeolocation").
			With("assoc_id", h.AssocID).
			With("device_id", deviceID).
			With("geolocation_id", geolocationID).
			With("latitude", lat).
			With("longitude", lng).
			With("altitude_meter", altMeter).
			With("accuracy_meter", accMeter).
			With("altitude_accuracy_meter", altAccMeter).
			With("heading", heading).
			With("speed_meter_per_sec", speedMeterPerSec).
			With("recorded_at", recordedAt).
			Wrapf(err, "failed to set geolocation to firestore")
	}

	slog.Info(
		"geolocation stored",
		"assoc_id", h.AssocID,
		"device_id", deviceID,
		"geolocation_id", geolocationID,
		"latitude", lat,
		"longitude", lng,
		"altitude_meter", altMeter,
		"accuracy_meter", accMeter,
		"altitude_accuracy_meter", altAccMeter,
		"heading", heading,
		"speed_meter_per_sec", speedMeterPerSec,
		"recorded_at", recordedAt,
	)

	return nil
}

func (h *Hub) FetchLatestGeolocationByDeviceID(
	ctx context.Context,
	deviceID string,
) (*Geolocation, error) {
	geolocationID := h.AssocID + "_" + deviceID

	loc, err := repoFirestore.FetchGeolocationByID(
		ctx,
		h.FirestoreClient,
		geolocationID,
	)
	if err != nil {
		return nil, oops.
			In("geolocationhub.FetchLatestGeolocationByDeviceID").
			With("assoc_id", h.AssocID).
			With("device_id", deviceID).
			With("geolocation_id", geolocationID).
			Wrapf(err, "failed to fetch geolocation by device id")
	}

	return &Geolocation{
		DeviceID:              deviceID,
		Latitude:              loc.Latitude,
		Longitude:             loc.Longitude,
		AltitudeMeter:         loc.AltitudeMeter,
		AccuracyMeter:         loc.AccuracyMeter,
		AltitudeAccuracyMeter: loc.AltitudeAccuracyMeter,
		Heading:               loc.Heading,
		SpeedMeterPerSec:      loc.SpeedMeterPerSec,
		RecordedAt:            loc.UpdatedAt,
	}, nil
}
