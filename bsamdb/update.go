package bsamdb

import (
	"database/sql"
	"fmt"
	"sailing-assist-mie-api/utils"
)

// Update updates values by selecting the primary_key.
func (db *DbInfo) Update(table string, primary_key string, id interface{}, fields []Field) (*sql.Row, error) {
	exist, err := db.IsExist(table, primary_key, id)
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
		primary_key,
		lastSqlSelector,
	)

	return db.DB.QueryRow(sql, values...), nil
}
