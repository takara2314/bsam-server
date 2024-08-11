package geolocationhub

import (
	"cloud.google.com/go/firestore"
)

type Hub struct {
	AssocID         string
	FirestoreClient *firestore.Client
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
