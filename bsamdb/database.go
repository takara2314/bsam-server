package bsamdb

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"sailing-assist-mie-api/utils"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/xo/dburl"
)

type DbInfo struct {
	DB *sql.DB
}

type Field struct {
	Column  string
	Value   interface{}
	Value2d []interface{}
	ToHash  bool
}

func Open() (DbInfo, error) {
	db, err := dburl.Open(os.Getenv("DATABASE_URL"))
	return DbInfo{DB: db}, err
}

func (db *DbInfo) Insert(table string, fields []Field) (*sql.Row, error) {
	sql := "INSERT INTO %s (%s) VALUES (%s)"

	columns := make([]string, len(fields))
	values := make([]interface{}, len(fields))

	for i, field := range fields {
		columns[i] = field.Column

		if !utils.IsNil(field.Value) {
			switch field.Value.(type) {
			case time.Time:
				values[i] = field.Value.(time.Time).Format("2006-01-02 15:04:05")
			default:
				values[i] = field.Value
			}
			continue
		}

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

		return nil, errors.New("must contain a value or a value slice")
	}

	sql_selectors := make([]string, len(fields))
	for i := 0; i < len(fields); i++ {
		sql_selectors[i] = "$" + strconv.Itoa(i+1)
	}
	for i, field := range fields {
		if field.ToHash {
			sql_selectors[i] = fmt.Sprintf(
				"digest($%d, 'sha3-256')",
				i+1,
			)

		} else {
			sql_selectors[i] = fmt.Sprintf(
				"$%d",
				i+1,
			)
		}
	}

	sql = fmt.Sprintf(
		sql,
		table,
		utils.StringSliceToString(columns),
		utils.StringSliceToString(sql_selectors),
	)

	return db.DB.QueryRow(sql, values...), nil
}
