package handler

import (
	"context"
	"log/slog"
	"time"

	"github.com/takara2314/bsam-server/internal/game/common"
	"github.com/takara2314/bsam-server/pkg/auth"
	"github.com/takara2314/bsam-server/pkg/devicelib"
	"github.com/takara2314/bsam-server/pkg/domain"
	"github.com/takara2314/bsam-server/pkg/geolocationlib"
	"github.com/takara2314/bsam-server/pkg/racehub"
)

type RaceHandler struct {
	racehub.UnimplementedHandler
}

// 認証メッセージを受信したときの処理
// 1. トークンが問題ないか検証
// 2. 内部の協会デバイスからの参加なら許可
// 3. デバイスIDが問題ないか検証
// 4. 選手ロールなら、ほしいマーク数が 1以上10以下 か検証
// 5. デバイスID、ロール、自分のマーク番号を登録
// 6. デバイス情報を記録
// 7. クライアントに認証完了メッセージを送信
// 8. 選手ロールなら、マークの位置情報を送信
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
	associationID, err := auth.ParseJWT(input.Token, common.Env.JWTSecretKey)
	if err != nil {
		slog.Warn(
			"failed to authenticate client",
			"client", c,
			"error", err,
			"input", input,
		)

		// クライアントに認証失敗した旨を送信
		if err := c.WriteAuthResult(
			false, racehub.AuthResultFailedAuthToken,
		); err != nil {
			slog.Error(
				"failed to write auth result",
				"client", c,
				"error", err,
			)
		}

		c.Hub.Unregister(c)
		return
	}

	// 外部の協会デバイスからの参加は現在許可しない
	if associationID != c.Hub.AssociationID {
		slog.Warn(
			"outside association client tried to connect",
			"client", c,
			"association_id", associationID,
			"input", input,
		)

		// クライアントに認証失敗した旨を送信
		if err := c.WriteAuthResult(
			false, racehub.AuthResultOutsideAssociation,
		); err != nil {
			slog.Error(
				"failed to write auth result",
				"client", c,
				"error", err,
			)
		}

		c.Hub.Unregister(c)
		return
	}

	// デバイスIDの検証を行う
	role, MarkNo, valid := domain.RetrieveRoleAndMarkNo(input.DeviceID)
	if !valid {
		slog.Warn(
			"invalid device_id",
			"client", c,
			"device_id", input.DeviceID,
			"input", input,
		)

		// クライアントに認証失敗した旨を送信
		if err := c.WriteAuthResult(
			false, racehub.AuthResultInvalidDeviceID,
		); err != nil {
			slog.Error(
				"failed to write auth result",
				"client", c,
				"error", err,
			)
		}

		c.Hub.Unregister(c)
		return
	}

	// 選手ロールなら、ほしいマーク数が 1以上10以下 か検証
	if role == domain.RoleAthlete {
		if input.WantMarkCounts < 1 || input.WantMarkCounts > 10 {
			slog.Warn(
				"invalid want_mark_no",
				"client", c,
				"want_mark_counts", input.WantMarkCounts,
				"input", input,
			)

			// クライアントに認証失敗した旨を送信
			if err := c.WriteAuthResult(
				false, racehub.AuthResultInvalidWantMarkCounts,
			); err != nil {
				slog.Error(
					"failed to write auth result",
					"client", c,
					"error", err,
				)
			}
		}
	}

	authedAt := time.Now()

	// デバイスID、ロール、自分のマーク番号を登録
	c.Hub.Mu.Lock()
	c.DeviceID = input.DeviceID
	c.Role = role
	c.MarkNo = MarkNo
	c.WantMarkCounts = input.WantMarkCounts
	c.Authed = true
	c.Hub.Mu.Unlock()

	slog.Info(
		"client authenticated",
		"client", c,
		"device_id", c.DeviceID,
		"role", c.Role,
		"mark_no", c.MarkNo,
		"input", input,
	)

	ctx := context.Background()

	// デバイス情報を記録
	if err := devicelib.StoreDevice(
		ctx,
		common.FirestoreClient,
		associationID,
		c.Hub.ID,
		c.DeviceID,
		c.ID,
		authedAt,
	); err != nil {
		slog.Error(
			"failed to store device",
			"client", c,
			"error", err,
			"racehub_id", c.Hub.ID,
			"device_id", c.DeviceID,
			"client_id", c.ID,
			"authed_at", authedAt,
		)
	}

	// クライアントに認証完了メッセージを送信
	if err := c.WriteAuthResult(
		true, racehub.AuthResultOK,
	); err != nil {
		slog.Error(
			"failed to write auth result",
			"client", c,
			"error", err,
		)
	}

	// 選手ロールなら、マークの位置情報を送信
	if c.Role == domain.RoleAthlete {
		time.Sleep(10 * time.Millisecond)
		if err := c.WriteMarkGeolocations(); err != nil {
			slog.Error(
				"failed to write mark geolocations",
				"client", c,
				"error", err,
			)
		}
	}
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

	ctx := context.Background()

	// 位置情報を記録
	if err := geolocationlib.StoreGeolocation(
		ctx,
		common.FirestoreClient,
		c.Hub.AssociationID,
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
