package event

import (
	"context"
	"log/slog"

	"github.com/takara2314/bsam-server/internal/game/common"
	"github.com/takara2314/bsam-server/pkg/devicelib"
	"github.com/takara2314/bsam-server/pkg/racehub"
)

type RaceEvent struct {
	racehub.UnimplementedEvent
}

// クライアントが接続したときの処理
func (r *RaceEvent) Register(c *racehub.Client) {
	slog.Info(
		"client registered",
		"client", c,
	)
}

// クライアントが切断するときの処理
func (r *RaceEvent) Unregister(c *racehub.Client) {
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
func (r *RaceEvent) ManageRaceStatusTaskReceived(
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
