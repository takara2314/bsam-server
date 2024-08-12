package devicehub

import (
	"context"
	"log/slog"
	"time"

	"github.com/samber/oops"
	repoFirestore "github.com/takara2314/bsam-server/pkg/infrastructure/repository/firestore"
)

type Device struct {
	HubID    string
	ClientID string
	AuthedAt time.Time
}

func (h *Hub) StoreDevice(
	ctx context.Context,
	raceHubID string,
	deviceID string,
	clientID string,
	authedAt time.Time,
) error {
	firestoreDeviceID := h.AssociationID + "_" + deviceID

	if err := repoFirestore.SetDevice(
		ctx,
		h.FirestoreClient,
		firestoreDeviceID,
		raceHubID,
		clientID,
		authedAt,
	); err != nil {
		return oops.
			In("geolocationhub.StoreDevice").
			With("association_id", h.AssociationID).
			With("device_id", deviceID).
			With("firestore_device_id", firestoreDeviceID).
			With("racehub_id", raceHubID).
			With("client_id", clientID).
			With("authed_at", authedAt).
			Wrapf(err, "failed to set device to firestore")
	}

	slog.Info(
		"device stored",
		"association_id", h.AssociationID,
		"device_id", deviceID,
		"firestore_device_id", firestoreDeviceID,
		"racehub_id", raceHubID,
		"client_id", clientID,
		"authed_at", authedAt,
	)

	return nil
}

func (h *Hub) FetchLatestDeviceByDeviceID(
	ctx context.Context,
	deviceID string,
) (*Device, error) {
	firestoreDeviceID := h.AssociationID + "_" + deviceID

	loc, err := repoFirestore.FetchDeviceByID(
		ctx,
		h.FirestoreClient,
		firestoreDeviceID,
	)
	if err != nil {
		return nil, oops.
			In("geolocationhub.FetchLatestDeviceByDeviceID").
			With("association_id", h.AssociationID).
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

func (h *Hub) DeleteFirestoreDeviceByDeviceID(
	ctx context.Context,
	deviceID string,
) error {
	firestoreDeviceID := h.AssociationID + "_" + deviceID

	if err := repoFirestore.DeleteDeviceByID(
		ctx,
		h.FirestoreClient,
		firestoreDeviceID,
	); err != nil {
		return oops.
			In("geolocationhub.DeleteFirestoreDeviceByDeviceID").
			With("association_id", h.AssociationID).
			With("device_id", deviceID).
			With("firestore_device_id", firestoreDeviceID).
			Wrapf(err, "failed to delete device by device id")
	}

	slog.Info(
		"device deleted",
		"association_id", h.AssociationID,
		"device_id", deviceID,
		"firestore_device_id", firestoreDeviceID,
	)

	return nil
}
