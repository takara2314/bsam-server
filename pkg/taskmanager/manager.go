package taskmanager

import (
	"context"
	"sync"

	"cloud.google.com/go/firestore"
)

const (
	PrefixManager     = "MANAGER_"
	PrefixAssociation = "ASSOCIATION_"
)

type Manager struct {
	ID               string
	AssociationID    string
	FirestoreClient  *firestore.Client
	Mu               sync.Mutex
	subscribeHandler SubscribeHandler
}

func NewManager(firestoreClient *firestore.Client) *Manager {
	return &Manager{
		ID:              "unknown",
		FirestoreClient: firestoreClient,
	}
}

func (m *Manager) StartManager(
	id string,
	associationID string,
	errorPipe chan error,
) {
	m.Mu.Lock()
	m.ID = id
	m.AssociationID = associationID
	m.Mu.Unlock()

	ctx := context.Background()
	go m.subscribeTasks(ctx, errorPipe)
}
