package firestore

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/samber/oops"
)

type Device struct {
	ID        string    `firestore:"-"`
	HubID     string    `firestore:"hubID"`
	ClientID  string    `firestore:"clientID"`
	UpdatedAt time.Time `firestore:"updatedAt"`
}

func SetDevice(
	ctx context.Context,
	client *firestore.Client,
	id string,
	hubID string,
	clientID string,
	updatedAt time.Time,
) error {
	_, err := client.Collection("devices").Doc(id).Set(ctx, Device{
		ID:        id,
		HubID:     hubID,
		ClientID:  clientID,
		UpdatedAt: updatedAt,
	})
	if err != nil {
		return oops.
			In("firestore.SetClient").
			Wrapf(err, "failed to set device")
	}

	return nil
}

func FetchDeviceByID(ctx context.Context, client *firestore.Client, id string) (*Device, error) {
	doc, err := client.Collection("devices").Doc(id).Get(ctx)
	if err != nil {
		return nil, oops.
			In("firestore.FetchDeviceByID").
			Wrapf(err, "failed to fetch device")
	}

	var c Device
	err = doc.DataTo(&c)
	if err != nil {
		return nil, oops.
			In("firestore.FetchDeviceByID").
			Wrapf(err, "failed to convert data to device")
	}

	c.ID = doc.Ref.ID

	return &c, err
}

func DeleteDeviceByID(ctx context.Context, client *firestore.Client, id string) error {
	_, err := client.Collection("devices").Doc(id).Delete(ctx)
	if err != nil {
		return oops.
			In("firestore.DeleteDeviceByID").
			Wrapf(err, "failed to delete device")
	}

	return nil
}
