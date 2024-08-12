package event

import (
	"context"
	"log/slog"

	"github.com/takara2314/bsam-server/internal/game/common"
	"github.com/takara2314/bsam-server/pkg/devicehub"
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
	deviceHub := devicehub.NewHub(
		c.Hub.AssociationID,
		common.FirestoreClient,
	)
	ctx := context.Background()

	// デバイス情報をFirestoreから削除
	// デバイスIDが存在しない場合等のエラーはスルー
	// (未認証のデバイスを削除する場合など)
	_ = deviceHub.DeleteFirestoreDeviceByDeviceID(
		ctx,
		c.DeviceID,
	)

	slog.Info(
		"client unregistered",
		"client", c,
	)
}
