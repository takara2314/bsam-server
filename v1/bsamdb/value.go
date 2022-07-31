package bsamdb

import (
	"bsam-server/utils"
	"fmt"
	"strconv"
	"time"
)

// valueToQueryable converts fields to queryable interface{} slice.
func valueToQueryable(fields []Field) ([]string, []interface{}, error) {
	columns := make([]string, len(fields))
	values := make([]interface{}, len(fields))

	// Convert to insertable form.
	for i, field := range fields {
		columns[i] = field.Column

		// The type of the value is not array.
		if !utils.IsNil(field.Value) {
			switch field.Value.(type) {
			case time.Time:
				values[i] = field.Value.(time.Time).Format("2006-01-02 15:04:05")
			default:
				values[i] = field.Value
			}
			continue
		}

		// The type of the value is array.
		//   {"a", "b", "c"} -> '{"a", "b", "c"}'
		if len(field.Value2d) != 0 {
			str := "{"
			last_loop := len(field.Value2d) - 1

			for j, item := range field.Value2d {
				switch item.(type) {
				case string:
					str += "\"" + item.(string) + "\""
				case time.Time:
					str += item.(time.Time).Format("2006-01-02 15:04:05")
				default:
					str += fmt.Sprint(item)
				}

				if last_loop != j {
					str += ", "
				} else {
					str += "}"
				}
			}

			values[i] = str
			continue
		}

		values[i] = "{}"
	}

	return columns, values, nil
}

// createSQLSelectors creates sql selectors when querying.
func createSQLSelector(fields []Field) []string {
	sqlSelectors := make([]string, len(fields))

	for i := 0; i < len(fields); i++ {
		sqlSelectors[i] = "$" + strconv.Itoa(i+1)
	}
	for i, field := range fields {
		// Insert after SHA3-256 hashed if $ToHash is true.
		if field.ToHash {
			sqlSelectors[i] = fmt.Sprintf(
				"digest($%d, 'sha3-256')",
				i+1,
			)

		} else {
			sqlSelectors[i] = fmt.Sprintf(
				"$%d",
				i+1,
			)
		}
	}

	return sqlSelectors
}
