package geolocationhub

import (
	"cloud.google.com/go/firestore"
)

type Hub struct {
	AssociationID   string
	FirestoreClient *firestore.Client
}

func NewHub(
	AssociationID string,
	firestoreClient *firestore.Client,
) *Hub {
	return &Hub{
		AssociationID:   AssociationID,
		FirestoreClient: firestoreClient,
	}
}
