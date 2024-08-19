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
	TaskTypeManageRaceStatus = "manage_race_start"
	TaskTypeManageNextMark   = "manage_next_mark"
)

type Hub struct {
	ID            string
	AssociationID string
	Clients       map[string]*Client
	Started       bool
	taskManager   *taskmanager.Manager
	clientEvent   ClientEvent
	serverEvent   ServerEvent
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
	clientEvent ClientEvent,
	serverEvent ServerEvent,
	handler Handler,
	action Action,
) *Hub {
	id := ulid.Make().String()

	slog.Info(
		"creating new hub",
		"id", id,
		"association_id", associationID,
	)

	// TODO: 協会のレース開始状態を取得する (from firestore)

	hub := &Hub{
		ID:            id,
		AssociationID: associationID,
		Clients:       make(map[string]*Client),
		taskManager:   tm,
		clientEvent:   clientEvent,
		serverEvent:   serverEvent,
		handler:       handler,
		action:        action,
	}

	tm.SetSubscribeHandler(hub.subscribeHandler)

	errorCh := make(chan error)
	go tm.StartManager(id, associationID, errorCh)

	return hub
}

type ServerEvent interface {
	ManageRaceStatusTaskReceived(*Hub, *ManageRaceStatusTaskMessage)
	ManageNextMarkTaskReceived(*Hub, *ManageNextMarkTaskMessage)
}

type UnimplementedServerEvent struct{}

type ManageRaceStatusTaskMessage struct {
	Started bool `json:"started"`
}

type ManageNextMarkTaskMessage struct {
	TargetDeviceID string `json:"target_device_id"`
	NextMarkNo     int    `json:"next_mark_no"`
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
		h.serverEvent.ManageRaceStatusTaskReceived(h, &ManageRaceStatusTaskMessage{})

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

func (h *Hub) PublishManageNextMarkTask(ctx context.Context, targetTaskManagerID string, targetDeviceID string, nextMarkNo int) error {
	payload, err := sonic.Marshal(&ManageNextMarkTaskMessage{
		TargetDeviceID: targetDeviceID,
		NextMarkNo:     nextMarkNo,
	})
	if err != nil {
		return err
	}

	if err := h.taskManager.PublishToManager(
		ctx,
		targetTaskManagerID,
		TaskTypeManageNextMark,
		payload,
	); err != nil {
		return err
	}

	return nil
}
