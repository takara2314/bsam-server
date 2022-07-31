package bsamdb

import (
	"database/sql"
	"errors"
	"os"

	_ "github.com/lib/pq"
	"github.com/xo/dburl"
)

var (
	ErrRecordNotFound = errors.New("this records is not found")
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

// Open connects to the database.
func Open() (DbInfo, error) {
	db, err := dburl.Open(os.Getenv("DATABASE_URL"))
	return DbInfo{DB: db}, err
}
