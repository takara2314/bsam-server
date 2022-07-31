package inspector

import (
	"os"

	"github.com/xo/dburl"

	_ "github.com/lib/pq"
)

// IsSameGroup returns that its group id and request token's group id is same.
func (ins *Inspector) IsSameGroup(groupID string) bool {
	if ins.Token.UserID.String == "" {
		return false
	}

	db, err := dburl.Open(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var tokenGroupID string
	row := db.QueryRow("SELECT group_id FROM users WHERE id = $1", ins.Token.UserID.String)
	err = row.Scan(&tokenGroupID)
	if err != nil {
		return false
	}

	if tokenGroupID != groupID {
		return false
	}

	return true
}
