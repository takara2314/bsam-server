package inspector

import (
	"os"

	"github.com/lib/pq"
	"github.com/xo/dburl"
)

// HasToken checks that its request header contains token and it is correct.
func (ins *Inspector) HasToken() string {
	const noAuthorization = "The token is required."
	const wrongAuthorization = "The authorization type is not supported."
	const wrongToken = "This token is wrong."

	auth := ins.Request.Header.Get("Authorization")

	if auth == "" {
		return noAuthorization
	}

	if len(auth) < 8 {
		return wrongAuthorization
	}

	if auth[:6] != "Bearer" {
		return wrongAuthorization
	}

	token := auth[7:]
	ins.Token.Token = token

	db, err := dburl.Open(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	row := db.QueryRow("SELECT * FROM tokens WHERE token = $1", token)
	err = row.Scan(&ins.Token.Token, pq.Array(&ins.Token.Permissions), &ins.Token.UserId, &ins.Token.Description)
	if err != nil {
		ins.IsTokenEnabled = false
		return wrongToken
	}

	ins.IsTokenEnabled = true
	return ""
}
