// create_table.go | for temporary

package main

import (
	"context"

	"cloud.google.com/go/bigquery"
)

func createTable(ctx context.Context, client *bigquery.Client, datasetID string, table *Table) error {
	metaData := &bigquery.TableMetadata{
		Schema: table.Schemes,
	}
	tableRef := client.Dataset(datasetID).Table(table.ID)
	if err := tableRef.Create(ctx, metaData); err != nil {
		return err
	}
	return nil
}
