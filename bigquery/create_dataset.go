// create_dataset.go | for temporary

package main

import (
	"context"

	"cloud.google.com/go/bigquery"
)

func createDataset(client *bigquery.Client, ctx context.Context, datasetID string, region string) error {
	meta := &bigquery.DatasetMetadata{
		Location: region,
	}
	if err := client.Dataset(datasetID).Create(ctx, meta); err != nil {
		return err
	}
	return nil
}
