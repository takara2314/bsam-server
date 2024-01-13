// main.go | for temporary

package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/option"
)

type Table struct {
	ID      string
	Schemes bigquery.Schema
}

var (
	projectID  = "bsam-app"
	datasetIDs = []string{
		"sensor_logs",
	}
	tables = []*Table{
		{
			ID: "location_logs",
			Schemes: bigquery.Schema{
				{Name: "user_id", Type: bigquery.StringFieldType, Required: false},
				{Name: "client_id", Type: bigquery.StringFieldType, Required: true},
				{Name: "role_id", Type: bigquery.IntegerFieldType, Required: true},
				{Name: "mark_no", Type: bigquery.IntegerFieldType, Required: false},
				{Name: "latitude", Type: bigquery.FloatFieldType, Required: true},
				{Name: "longitude", Type: bigquery.FloatFieldType, Required: true},
				{Name: "accuracy", Type: bigquery.FloatFieldType, Required: true},
				{Name: "heading", Type: bigquery.FloatFieldType, Required: true},
				{Name: "recorded_at", Type: bigquery.DateTimeFieldType, Required: true},
			},
		},
	}
	region = "asia-northeast1"
)

func main() {
	ctx := context.Background()
	auth := option.WithCredentialsFile("./bsam-app-d23e7e5025e7.json")

	client, err := bigquery.NewClient(ctx, projectID, auth)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	for _, datasetID := range datasetIDs {
		err := createDataset(client, ctx, datasetID, region)
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, table := range tables {
			err := createTable(client, ctx, datasetID, table, region)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}
