package bsamdb

import (
	"bsam-server/utils"
	"database/sql"
	"fmt"
)

// Update updates values by selecting the primaryKey.
func (db *DbInfo) Update(table string, primaryKey string, id interface{}, fields []Field) (*sql.Row, error) {
	exist, err := db.IsExist(table, primaryKey, id)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, ErrRecordNotFound
	}

	sql := "UPDATE %s SET %s WHERE %s = %s"

	// Convert to insertable form.
	columns, values, err := valueToQueryable(fields)
	if err != nil {
		return nil, err
	}
	values = append(values, id)

	// Create sql selectors.
	sqlSelectors := createSQLSelector(fields)
	lastSqlSelector := fmt.Sprintf("$%d", len(sqlSelectors)+1)

	// Create the alter statement.
	alter, err := utils.CreateStrSliceEqualStrSlice(columns, sqlSelectors)
	if err != nil {
		return nil, err
	}

	// Execute sql query.
	sql = fmt.Sprintf(
		sql,
		table,
		alter,
		primaryKey,
		lastSqlSelector,
	)

	return db.DB.QueryRow(sql, values...), nil
}

// Update updates all values.
func (db *DbInfo) UpdateAll(table string, fields []Field) (*sql.Row, error) {
	sql := "UPDATE %s SET %s"

	// Convert to insertable form.
	columns, _, err := valueToQueryable(fields)
	if err != nil {
		return nil, err
	}

	// Create sql selectors.
	sqlSelectors := createSQLSelector(fields)

	// Create the alter statement.
	alter, err := utils.CreateStrSliceEqualStrSlice(columns, sqlSelectors)
	if err != nil {
		return nil, err
	}

	// Execute sql query.
	sql = fmt.Sprintf(
		sql,
		table,
		alter,
	)

	return db.DB.QueryRow(sql), nil
}
