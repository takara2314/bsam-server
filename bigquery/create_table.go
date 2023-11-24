// create_table.go | for temporary

package main

import (
	"context"

	"cloud.google.com/go/bigquery"
)

func createTable(client *bigquery.Client, ctx context.Context, datasetID string, table *Table, region string) error {
	metaData := &bigquery.TableMetadata{
		Schema: table.Schemes,
	}
	tableRef := client.Dataset(datasetID).Table(table.ID)
	if err := tableRef.Create(ctx, metaData); err != nil {
		return err
	}
	return nil
}
