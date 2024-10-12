package devicelib

import (
	"context"
	"log/slog"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/samber/oops"
	repoFirestore "github.com/takara2314/bsam-server/pkg/infrastructure/repository/firestore"
	"github.com/takara2314/bsam-server/pkg/racelib"
)

type Device struct {
	HubID    string
	ClientID string
	AuthedAt time.Time
}

func StoreDevice(
	ctx context.Context,
	firestoreClient *firestore.Client,
	associationID string,
	raceHubID string,
	deviceID string,
	clientID string,
	authedAt time.Time,
) error {
	firestoreDeviceID := associationID + "_" + deviceID

	// デバイスコレクションにデバイス情報を設定
	if err := repoFirestore.SetDevice(
		ctx,
		firestoreClient,
		firestoreDeviceID,
		raceHubID,
		clientID,
		authedAt,
	); err != nil {
		return oops.
			In("geolocationlib.StoreDevice").
			With("association_id", associationID).
			With("device_id", deviceID).
			With("firestore_device_id", firestoreDeviceID).
			With("racehub_id", raceHubID).
			With("client_id", clientID).
			With("authed_at", authedAt).
			Wrapf(err, "failed to set device to firestore")
	}

	// レースコレクションにデバイスIDを設定
	if err := racelib.AddDeviceIDToRace(
		ctx,
		firestoreClient,
		associationID,
		deviceID,
	); err != nil {
		return oops.
			In("geolocationlib.StoreDevice").
			With("association_id", associationID).
			With("device_id", deviceID).
			Wrapf(err, "failed to add device id to race")
	}

	slog.Info(
		"device stored",
		"association_id", associationID,
		"device_id", deviceID,
		"firestore_device_id", firestoreDeviceID,
		"racehub_id", raceHubID,
		"client_id", clientID,
		"authed_at", authedAt,
	)

	return nil
}

func FetchLatestDeviceByDeviceID(
	ctx context.Context,
	firestoreClient *firestore.Client,
	associationID string,
	deviceID string,
) (*Device, error) {
	firestoreDeviceID := associationID + "_" + deviceID

	loc, err := repoFirestore.FetchDeviceByID(
		ctx,
		firestoreClient,
		firestoreDeviceID,
	)
	if err != nil {
		return nil, oops.
			In("geolocationlib.FetchLatestDeviceByDeviceID").
			With("association_id", associationID).
			With("device_id", deviceID).
			With("firestore_device_id", firestoreDeviceID).
			Wrapf(err, "failed to fetch device by device id")
	}

	return &Device{
		HubID:    loc.HubID,
		ClientID: loc.ClientID,
		AuthedAt: loc.UpdatedAt,
	}, nil
}

func DeleteFirestoreDeviceByDeviceID(
	ctx context.Context,
	firestoreClient *firestore.Client,
	associationID string,
	deviceID string,
) error {
	firestoreDeviceID := associationID + "_" + deviceID

	if err := repoFirestore.DeleteDeviceByID(
		ctx,
		firestoreClient,
		firestoreDeviceID,
	); err != nil {
		return oops.
			In("geolocationlib.DeleteFirestoreDeviceByDeviceID").
			With("association_id", associationID).
			With("device_id", deviceID).
			With("firestore_device_id", firestoreDeviceID).
			Wrapf(err, "failed to delete device by device id")
	}

	if err := racelib.RemoveDeviceIDFromRace(
		ctx,
		firestoreClient,
		associationID,
		deviceID,
	); err != nil {
		return oops.
			In("geolocationlib.DeleteFirestoreDeviceByDeviceID").
			With("association_id", associationID).
			With("device_id", deviceID).
			Wrapf(err, "failed to remove device id from race")
	}

	slog.Info(
		"device deleted",
		"association_id", associationID,
		"device_id", deviceID,
		"firestore_device_id", firestoreDeviceID,
	)

	return nil
}
