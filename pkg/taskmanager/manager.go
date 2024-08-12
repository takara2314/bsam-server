package taskmanager

import (
	"context"
	"sync"

	"cloud.google.com/go/firestore"
)

type Manager struct {
	ID               string
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

func (m *Manager) StartManager(id string, errorPipe chan error) {
	m.Mu.Lock()
	m.ID = id
	m.Mu.Unlock()

	ctx := context.Background()
	go m.subscribeTasks(ctx, errorPipe)
}
