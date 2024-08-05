package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/takara2314/bsam-server/internal/api/common"
	"github.com/takara2314/bsam-server/pkg/infrastructure/repository/firestore"
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

func GeolocationPOST(c *gin.Context, assocID string, req GeolocationPOSTRequest) {
	geolocationID := assocID + "_" + req.DeviceID

	if err := firestore.SetGeolocation(
		c,
		common.FirestoreClient,
		geolocationID,
		req.Latitude,
		req.Longitude,
		req.AltitudeMeter,
		req.AccuracyMeter,
		req.AltitudeAccuracyMeter,
		req.Heading,
		req.SpeedMeterPerSec,
		req.RecordedAt,
	); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "failed to set geolocation to firestore",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "OK",
		"geolocation": req,
	})
}
