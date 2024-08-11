package handler

import (
	"context"
	"log/slog"

	"github.com/takara2314/bsam-server/internal/game/common"
	"github.com/takara2314/bsam-server/pkg/auth"
	"github.com/takara2314/bsam-server/pkg/domain"
	"github.com/takara2314/bsam-server/pkg/geolocationhub"
	"github.com/takara2314/bsam-server/pkg/racehub"
)

type RaceHandler struct {
	racehub.UnimplementedHandler
}

// 認証メッセージを受信したときの処理
// 1. トークンが問題ないか検証
// 2. 内部の協会デバイスからの参加なら許可
// 3. デバイスID、ロール、自分のマーク番号を登録
// 4. クライアントに認証完了メッセージを送信
func (r *RaceHandler) Auth(
	c *racehub.Client,
	input *racehub.AuthInput,
) {
	slog.Info(
		"received auth message",
		"client", c,
		"input", input,
	)

	// JWTトークンを検証
	assocID, err := auth.ParseJWT(input.Token, common.Env.JWTSecretKey)
	if err != nil {
		slog.Warn(
			"failed to authenticate client",
			"client", c,
			"error", err,
			"input", input,
		)

		// クライアントに認証失敗した旨を送信
		c.SendAuthResult(
			false,
			c.DeviceID,
			c.Role,
			c.MyMarkNo,
			racehub.AuthResultFailedAuthToken,
		)

		c.Hub.Unregister(c)
		return
	}

	// 外部の協会デバイスからの参加は現在許可しない
	if assocID != c.Hub.AssocID {
		slog.Warn(
			"outside assoc client tried to connect",
			"client", c,
			"assoc_id", assocID,
			"input", input,
		)

		// クライアントに認証失敗した旨を送信
		c.SendAuthResult(
			false,
			c.DeviceID,
			c.Role,
			c.MyMarkNo,
			racehub.AuthResultOutsideAssoc,
		)

		c.Hub.Unregister(c)
		return
	}

	// デバイスIDの検証を行い、デバイスID、ロール、自分のマーク番号を登録
	c.Hub.Mu.Lock()
	var valid bool
	c.DeviceID = input.DeviceID

	c.Role, c.MyMarkNo, valid = domain.RetrieveRoleAndMyMarkNo(input.DeviceID)
	if !valid {
		slog.Warn(
			"invalid device_id",
			"client", c,
			"device_id", input.DeviceID,
			"input", input,
		)

		// クライアントに認証失敗した旨を送信
		c.SendAuthResult(
			false,
			c.DeviceID,
			c.Role,
			c.MyMarkNo,
			racehub.AuthResultInvalidDeviceID,
		)

		c.Hub.Unregister(c)
		return
	}
	c.Hub.Mu.Unlock()

	slog.Info(
		"client authenticated",
		"client", c,
		"device_id", c.DeviceID,
		"role", c.Role,
		"my_mark_no", c.MyMarkNo,
		"input", input,
	)

	// クライアントに認証完了メッセージを送信
	c.SendAuthResult(
		true,
		c.DeviceID,
		c.Role,
		c.MyMarkNo,
		racehub.AuthResultOK,
	)
}

// 位置情報を受信したときの処理
func (r *RaceHandler) PostGeolocation(
	c *racehub.Client,
	input *racehub.PostGeolocationInput,
) {
	slog.Info(
		"received post_geolocation message",
		"client", c,
		"input", input,
	)

	geoHub := geolocationhub.NewHub(
		c.Hub.AssocID,
		common.FirestoreClient,
	)
	ctx := context.Background()

	// 位置情報を記録
	if err := geoHub.StoreGeolocation(
		ctx,
		c.DeviceID,
		input.Latitude,
		input.Longitude,
		input.AltitudeMeter,
		input.AccuracyMeter,
		input.AltitudeAccuracyMeter,
		input.Heading,
		input.SpeedMeterPerSec,
		input.RecordedAt,
	); err != nil {
		slog.Warn(
			"failed to store geolocation",
			"client", c,
			"error", err,
			"input", input,
		)
		return
	}

	slog.Info(
		"geolocation saved",
		"client", c,
		"input", input,
	)
}
