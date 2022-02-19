package bsamdb

import (
	"fmt"
)

// IsExist checks the primary_key recorded its id is existed.
func (db *DbInfo) IsExist(table string, primary_key string, id interface{}) (bool, error) {
	sql := fmt.Sprintf(
		"SELECT COUNT(*) FROM %s WHERE %s = $1",
		table,
		primary_key,
	)

	row := db.DB.QueryRow(sql, id)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}
	return false, nil
}
