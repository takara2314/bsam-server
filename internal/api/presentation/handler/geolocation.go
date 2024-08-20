package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/takara2314/bsam-server/internal/api/common"
	"github.com/takara2314/bsam-server/pkg/geolocationlib"
)

type GeolocationPOSTRequest struct {
	DeviceID              string    `json:"device_id"`
	Latitude              float64   `json:"latitude"`
	Longitude             float64   `json:"longitude"`
	AltitudeMeter         float64   `json:"altitude_meter"`
	AccuracyMeter         float64   `json:"accuracy_meter"`
	AltitudeAccuracyMeter float64   `json:"altitude_accuracy_meter"`
	Heading               float64   `json:"heading"`
	SpeedMeterPerSec      float64   `json:"speed_meter_per_sec"`
	RecordedAt            time.Time `json:"recorded_at"`
}

func GeolocationPOST(c *gin.Context, associationID string, req GeolocationPOSTRequest) {
	// 位置情報を記録
	if err := geolocationlib.StoreGeolocation(
		c,
		common.FirestoreClient,
		associationID,
		req.DeviceID,
		req.Latitude,
		req.Longitude,
		req.AltitudeMeter,
		req.AccuracyMeter,
		req.AltitudeAccuracyMeter,
		req.Heading,
		req.SpeedMeterPerSec,
		req.RecordedAt,
	); err != nil {
		slog.Warn(
			"failed to store geolocation",
			"error", err,
			"request", req,
		)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "failed to post geolocation",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Created",
		"geolocation": req,
	})
}
