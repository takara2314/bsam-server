package racehub

import (
	"log/slog"
	"sync"

	"github.com/oklog/ulid/v2"
)

type Hub struct {
	ID            string
	AssociationID string
	clients       map[string]*Client
	handler       Handler
	action        Action
	Mu            sync.RWMutex
}

func (h *Hub) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("id", h.ID),
		slog.String("association_id", h.AssociationID),
		slog.Int("client_counts", len(h.clients)),
	)
}

func NewHub(associationID string, handler Handler, action Action) *Hub {
	id := ulid.Make().String()

	slog.Info(
		"creating new hub",
		"id", id,
		"association_id", associationID,
	)

	return &Hub{
		ID:            id,
		AssociationID: associationID,
		clients:       make(map[string]*Client),
		handler:       handler,
		action:        action,
	}
}
