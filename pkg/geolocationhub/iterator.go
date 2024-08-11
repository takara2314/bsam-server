package geolocationhub

import (
	"context"
	"strings"

	repoFirestore "github.com/takara2314/bsam-server/pkg/infrastructure/repository/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StreamIterator struct {
	hub   *GeolocationHub
	queue []Geolocation
	err   error
}

func (h *GeolocationHub) Snapshots(ctx context.Context) *StreamIterator {
	it := &StreamIterator{hub: h}
	go it.watchGeolocationCollection(ctx)
	return it
}

func (it *StreamIterator) watchGeolocationCollection(ctx context.Context) {
	itFirestore := it.hub.FirestoreClient.Collection("geolocations").Snapshots(ctx)
	for {
		snap, err := itFirestore.Next()

		if status.Code(err) == codes.DeadlineExceeded {
			return
		}
		if err != nil {
			it.err = err
			return
		}

		for _, change := range snap.Changes {
			var loc repoFirestore.Geolocation

			if !strings.HasPrefix(change.Doc.Ref.ID, it.hub.AssocID) {
				continue
			}

			err = change.Doc.DataTo(&loc)
			if err != nil {
				it.err = err
				return
			}

			it.queue = append(it.queue, Geolocation{
				DeviceID:              change.Doc.Ref.ID,
				Latitude:              loc.Latitude,
				Longitude:             loc.Longitude,
				AltitudeMeter:         loc.AltitudeMeter,
				AccuracyMeter:         loc.AccuracyMeter,
				AltitudeAccuracyMeter: loc.AltitudeAccuracyMeter,
				Heading:               loc.Heading,
				SpeedMeterPerSec:      loc.SpeedMeterPerSec,
				RecordedAt:            loc.UpdatedAt,
			})
		}
	}
}

func (i *StreamIterator) Next() (*Geolocation, error) {
	for {
		if i.err != nil {
			return nil, i.err
		}
		if len(i.queue) > 0 {
			geolocation := i.queue[0]
			i.queue = i.queue[1:]
			return &geolocation, nil
		}
	}
}
