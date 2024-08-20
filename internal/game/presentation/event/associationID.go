package event

import (
	"context"
	"log/slog"

	"github.com/takara2314/bsam-server/internal/game/common"
	"github.com/takara2314/bsam-server/pkg/devicelib"
	"github.com/takara2314/bsam-server/pkg/racehub"
)

type RaceClientEvent struct {
	racehub.UnimplementedClientEvent
}

type RaceServerEvent struct {
	racehub.UnimplementedServerEvent
}

// クライアントが接続したときの処理
func (r *RaceClientEvent) Register(c *racehub.Client) {
	slog.Info(
		"client registered",
		"client", c,
	)
}

// クライアントが切断するときの処理
func (r *RaceClientEvent) Unregister(c *racehub.Client) {
	ctx := context.Background()

	// デバイス情報をFirestoreから削除
	// デバイスIDが存在しない場合等のエラーはスルー
	// (未認証のデバイスを削除する場合など)
	_ = devicelib.DeleteFirestoreDeviceByDeviceID(
		ctx,
		common.FirestoreClient,
		c.Hub.AssociationID,
		c.DeviceID,
	)

	slog.Info(
		"client unregistered",
		"client", c,
	)
}

// レースの状態を管理するタスクを受信したときの処理
// 認証済み全員にレース開始アクションを送信する
func (r *RaceServerEvent) ManageRaceStatusTaskReceived(
	h *racehub.Hub,
	msg *racehub.ManageRaceStatusTaskMessage,
) {
	h.Mu.Lock()
	h.Started = msg.Started
	h.Mu.Unlock()

	for _, c := range h.Clients {
		if !c.Authed {
			continue
		}
		_ = c.WriteManageRaceStatus(msg.Started)
	}
}

// 次のマークの管理タスクを受信したときの処理
// 指定のデバイスに次のマークの情報を送信する
func (r *RaceServerEvent) ManageNextMarkTaskReceived(
	h *racehub.Hub,
	msg *racehub.ManageNextMarkTaskMessage,
) {
	for _, c := range h.Clients {
		if c.DeviceID == msg.TargetDeviceID {
			_ = c.WriteManageNextMark(msg.NextMarkNo)
			return
		}
	}

	slog.Error(
		"client not found",
		"target_device_id", msg.TargetDeviceID,
		"hub", h,
	)
}
