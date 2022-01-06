package inspector

import (
	"os"

	"github.com/xo/dburl"

	_ "github.com/lib/pq"
)

// IsSameGroup returns that its group id and request token's group id is same.
func (ins *Inspector) IsSameGroup(groupId string) bool {
	if ins.Token.UserId.String == "" {
		return false
	}

	db, err := dburl.Open(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var tokenGroupId string
	row := db.QueryRow("SELECT group_id FROM users WHERE id = $1", ins.Token.UserId.String)
	err = row.Scan(&tokenGroupId)
	if err != nil {
		return false
	}

	if tokenGroupId != groupId {
		return false
	}

	return true
}
