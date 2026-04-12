package racing

import (
	"context"
	"errors"
	"os"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"google.golang.org/api/option"
)

const (
	gcpProjectID = "bsam-app"
	// #nosec G101 -- This is a local fallback path for development only.
	localCredentialPath = "./bsam-app-d23e7e5025e7.json"
)

type BigQueryLogger struct {
	Client *bigquery.Client
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

	client, err := newBigQueryClient(ctx)
	if err != nil {
		panic(err)
	}

	return &BigQueryLogger{
		Client: client,
	}
}

func newBigQueryClient(ctx context.Context) (*bigquery.Client, error) {
	credentialPath := localCredentialPath
	if envCredentialPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); envCredentialPath != "" {
		credentialPath = envCredentialPath
	}

	// #nosec G703 -- Credential files are intentionally sourced from ADC env or the local dev fallback.
	_, err := os.Stat(credentialPath)
	if err == nil {
		// 認証ファイルがあるときは認証ファイルを使用 (GCP環境ではないとき)
		auth := option.WithAuthCredentialsFile(option.ServiceAccount, credentialPath)

		return bigquery.NewClient(ctx, gcpProjectID, auth)
	}

	if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	return bigquery.NewClient(ctx, gcpProjectID)
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

	err := inserter.Put(context.Background(), data)
	if err != nil {
		return err
	}

	return nil
}
