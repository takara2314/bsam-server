package inspector

import (
	"database/sql"
	"net/http"
)

type Inspector struct {
	Request        *http.Request
	IsTokenEnabled bool
	Token          struct {
		Token       string
		Type        string
		Permissions []string
		UserId      sql.NullString
		Description sql.NullString
	}
}
