package racehub

import (
	"log/slog"
	"sync"

	"github.com/oklog/ulid/v2"
	"github.com/takara2314/bsam-server/pkg/taskmanager"
)

type Hub struct {
	ID            string
	AssociationID string
	clients       map[string]*Client
	TaskManager   *taskmanager.Manager
	event         Event
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

func NewHub(
	associationID string,
	tm *taskmanager.Manager,
	event Event,
	handler Handler,
	action Action,
) *Hub {
	id := ulid.Make().String()

	slog.Info(
		"creating new hub",
		"id", id,
		"association_id", associationID,
	)

	hub := &Hub{
		ID:            id,
		AssociationID: associationID,
		clients:       make(map[string]*Client),
		TaskManager:   tm,
		event:         event,
		handler:       handler,
		action:        action,
	}

	tm.SetSubscribeHandler(hub.subscribeHandler)

	errorPipe := make(chan error)
	go tm.StartManager(id, errorPipe)

	return hub
}

func (h *Hub) subscribeHandler(msgType string, payload []byte) error {
	slog.Info(
		"received task",
		"hub", h,
		"msg_type", msgType,
		"payload", string(payload),
	)
	return nil
}
