package bsamdb

import (
	"fmt"
)

// IsExist checks the primary_key recorded its id is existed.
func (db *DbInfo) IsExist(table string, primaryKey string, id interface{}) (bool, error) {
	sql := fmt.Sprintf(
		"SELECT COUNT(*) FROM %s WHERE %s = $1",
		table,
		primaryKey,
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

// IsExist checks the primary_key with secondary_key recorded its id is existed.
func (db *DbInfo) IsExist2(table string, primaryKey string, id interface{}, secondaryKey string, value interface{}) (bool, error) {
	sql := fmt.Sprintf(
		"SELECT COUNT(*) FROM %s WHERE %s = $1 AND %s = $2",
		table,
		primaryKey,
		secondaryKey,
	)

	row := db.DB.QueryRow(sql, id, value)
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

// IsExistNotIt checks the primary_key recorded its id excepting $notKey value is existed.
func (db *DbInfo) IsExistNotIt(table string, primaryKey string, id interface{}, notKey string, value interface{}) (bool, error) {
	sql := fmt.Sprintf(
		"SELECT COUNT(*) FROM %s WHERE %s = $1 AND %s != $2",
		table,
		primaryKey,
		notKey,
	)

	row := db.DB.QueryRow(sql, id, value)
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
