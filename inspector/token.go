package inspector

import (
	"os"

	"sailing-assist-mie-api/message"

	"github.com/lib/pq"
	"github.com/xo/dburl"
)

// HasToken checks that its request header contains token and it is correct.
func (ins *Inspector) HasToken() string {
	auth := ins.Request.Header.Get("Authorization")

	if auth == "" {
		return message.TokenNotFound
	}

	if len(auth) < 8 {
		return message.AuthorizationTypeInvalid
	}

	if auth[:6] != "Bearer" {
		return message.AuthorizationTypeInvalid
	}

	token := auth[7:]
	ins.Token.Token = token

	db, err := dburl.Open(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	row := db.QueryRow("SELECT * FROM tokens WHERE token = $1", token)
	err = row.Scan(&ins.Token.Token, &ins.Token.Type, pq.Array(&ins.Token.Permissions), &ins.Token.UserId, &ins.Token.Description)
	if err != nil {
		ins.IsTokenEnabled = false
		return message.WrongToken
	}

	ins.IsTokenEnabled = true
	return ""
}
