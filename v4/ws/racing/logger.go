package racing

import (
	"context"
	"os"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"google.golang.org/api/option"
)

const (
	gcpProjectID = "bsam-app"
)

type BigQueryLogger struct {
	Client *bigquery.Client
	Ctx    context.Context
}

type LocationLogsDAO struct {
	UserID     bigquery.NullString `bigquery:"user_id"`
	ClientID   string              `bigquery:"client_id"`
	RoleID     int64               `bigquery:"role_id"`
	MarkNo     bigquery.NullInt64  `bigquery:"mark_no"`
	Latitude   float64             `bigquery:"latitude"`
	Longitude  float64             `bigquery:"longitude"`
	Accuracy   float64             `bigquery:"accuracy"`
	Heading    float64             `bigquery:"heading"`
	RecordedAt civil.DateTime      `bigquery:"recorded_at"`
}

func NewBigQueryLogger() *BigQueryLogger {
	ctx := context.Background()
	var client *bigquery.Client

	if _, err := os.Stat("./bsam-app-d23e7e5025e7.json"); !os.IsNotExist(err) {
		// 認証ファイルがあるときは認証ファイルを使用 (GCP環境ではないとき)
		auth := option.WithCredentialsFile("./bsam-app-d23e7e5025e7.json")
		client, err = bigquery.NewClient(ctx, gcpProjectID, auth)
		if err != nil {
			panic(err)
		}

	} else {
		client, err = bigquery.NewClient(ctx, gcpProjectID)
		if err != nil {
			panic(err)
		}
	}

	return &BigQueryLogger{
		Client: client,
		Ctx:    ctx,
	}
}

func (l *BigQueryLogger) logLocation(c *Client) error {
	tableRef := l.Client.Dataset("sensor_logs").Table("location_logs")

	userID := bigquery.NullString{StringVal: c.UserID, Valid: true}
	if c.UserID == "" {
		userID.Valid = false
	}

	markNo := bigquery.NullInt64{Int64: int64(c.MarkNo), Valid: true}
	if c.MarkNo == -1 {
		markNo.Valid = false
	}

	data := LocationLogsDAO{
		UserID:     userID,
		ClientID:   c.ID,
		RoleID:     int64(c.getRoleID()),
		MarkNo:     markNo,
		Latitude:   c.Location.Lat,
		Longitude:  c.Location.Lng,
		Accuracy:   c.Location.Acc,
		Heading:    c.Location.Heading,
		RecordedAt: civil.DateTimeOf(time.Now()),
	}

	inserter := tableRef.Inserter()
	if err := inserter.Put(l.Ctx, data); err != nil {
		return err
	}

	return nil
}
