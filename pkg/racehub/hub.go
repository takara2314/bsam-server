package racehub

import (
	"context"
	"log/slog"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/oklog/ulid/v2"
	"github.com/takara2314/bsam-server/pkg/taskmanager"
)

const (
	TaskTypeManageRaceStatus = "race_start"
)

type Hub struct {
	ID            string
	AssociationID string
	Clients       map[string]*Client
	Started       bool
	taskManager   *taskmanager.Manager
	event         Event
	handler       Handler
	action        Action
	Mu            sync.RWMutex
}

func (h *Hub) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("id", h.ID),
		slog.String("association_id", h.AssociationID),
		slog.Int("client_counts", len(h.Clients)),
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
		Clients:       make(map[string]*Client),
		taskManager:   tm,
		event:         event,
		handler:       handler,
		action:        action,
	}

	tm.SetSubscribeHandler(hub.subscribeHandler)

	errorCh := make(chan error)
	go tm.StartManager(id, associationID, errorCh)

	return hub
}

type ManageRaceStatusTaskMessage struct {
	Started bool `json:"started"`
}

func (h *Hub) subscribeHandler(taskType string, payload []byte) error {
	slog.Info(
		"received task",
		"hub", h,
		"task_type", taskType,
		"payload", string(payload),
	)

	switch taskType {
	case TaskTypeManageRaceStatus:
		h.event.ManageRaceStatusTaskReceived(h, &ManageRaceStatusTaskMessage{})

	default:
		slog.Warn(
			"unsupported task type",
			"hub", h,
			"task_type", taskType,
			"payload", string(payload),
		)
	}

	return nil
}

func (h *Hub) PublishManageRaceStatusTask(ctx context.Context, started bool) error {
	payload, err := sonic.Marshal(&ManageRaceStatusTaskMessage{
		Started: started,
	})
	if err != nil {
		return err
	}

	if err := h.taskManager.PublishToAssociation(
		ctx,
		h.AssociationID,
		TaskTypeManageRaceStatus,
		payload,
	); err != nil {
		return err
	}

	return nil
}
