package bsamdb

import (
	"database/sql"
	"fmt"
	"sailing-assist-mie-api/utils"
)

// Insert runs a insert query.
func (db *DbInfo) Insert(table string, fields []Field) (*sql.Row, error) {
	sql := "INSERT INTO %s (%s) VALUES (%s)"

	// Convert to insertable form.
	columns, values, err := valueToQueryable(fields)
	if err != nil {
		return nil, err
	}

	// Create sql selectors.
	sqlSelectors := createSQLSelector(fields)

	// Execute sql query.
	sql = fmt.Sprintf(
		sql,
		table,
		utils.StringSliceToString(columns),
		utils.StringSliceToString(sqlSelectors),
	)

	fmt.Println(sql)
	fmt.Println(values...)

	return db.DB.QueryRow(sql, values...), nil
}
