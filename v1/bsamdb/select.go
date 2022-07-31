package bsamdb

import (
	"bsam-server/utils"
	"database/sql"
	"fmt"
)

// Select runs a insert query.
func (db *DbInfo) Select(table string, fields []Field) (*sql.Rows, error) {
	sql := "SELECT * FROM %s"
	sqlWithWhere := "SELECT * FROM %s WHERE %s"

	// Convert to insertable form.
	columns, values, err := valueToQueryable(fields)
	if err != nil {
		return nil, err
	}

	// Create sql selectors.
	sqlSelectors := createSQLSelector(fields)

	// If no where selectors.
	if len(fields) == 0 {
		// Execute sql query.
		sql = fmt.Sprintf(
			sql,
			table,
		)

		return db.DB.Query(sql, values...)
	}

	// Create the alter statement.
	alter, err := utils.CreateStrSliceEqualStrSlice(columns, sqlSelectors)
	if err != nil {
		return nil, err
	}

	// Execute sql query.
	sqlWithWhere = fmt.Sprintf(
		sqlWithWhere,
		table,
		alter,
	)

	return db.DB.Query(sqlWithWhere, values...)
}

// SelectSpecified runs a insert query with specify.
func (db *DbInfo) SelectSpecified(table string, fields []Field, columnsSpecified []string) (*sql.Rows, error) {
	columnsSpecifiedStr := utils.StringSliceToString(columnsSpecified)

	sql := "SELECT " + columnsSpecifiedStr + " FROM %s"
	sqlWithWhere := "SELECT " + columnsSpecifiedStr + " FROM %s WHERE %s"

	// Convert to insertable form.
	columns, values, err := valueToQueryable(fields)
	if err != nil {
		return nil, err
	}

	// Create sql selectors.
	sqlSelectors := createSQLSelector(fields)

	// If no where selectors.
	if len(fields) == 0 {
		// Execute sql query.
		sql = fmt.Sprintf(
			sql,
			table,
		)

		return db.DB.Query(sql, values...)
	}

	// Create the alter statement.
	alter, err := utils.CreateStrSliceEqualStrSlice(columns, sqlSelectors)
	if err != nil {
		return nil, err
	}

	// Execute sql query.
	sqlWithWhere = fmt.Sprintf(
		sqlWithWhere,
		table,
		alter,
	)

	return db.DB.Query(sqlWithWhere, values...)
}
