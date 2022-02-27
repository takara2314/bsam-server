package inspector

import (
	"os"

	"github.com/lib/pq"
	"github.com/xo/dburl"
)

// HasToken checks that its request header contains token and it is correct.
func (ins *Inspector) HasToken() (bool, error) {
	auth := ins.Request.Header.Get("Authorization")

	if auth == "" {
		return false, ErrTokenNotFound
	}

	if len(auth) < 8 {
		return false, ErrAuthorizationTypeInvalid
	}

	if auth[:6] != "Bearer" {
		return false, ErrAuthorizationTypeInvalid
	}

	token := auth[7:]
	ins.Token.Token = token

	db, err := dburl.Open(os.Getenv("DATABASE_URL"))
	if err != nil {
		return false, err
	}
	defer db.Close()

	row := db.QueryRow("SELECT * FROM tokens WHERE token = $1", token)
	err = row.Scan(&ins.Token.Token, &ins.Token.Type, pq.Array(&ins.Token.Permissions), &ins.Token.UserId, &ins.Token.Description)
	if err != nil {
		ins.IsTokenEnabled = false
		return false, ErrWrongToken
	}

	ins.IsTokenEnabled = true
	return true, nil
}

// FetchGroupId fetches groupId identified its request header.
func (ins *Inspector) FetchGroupId() error {
	_, err := ins.HasToken()
	if err != nil {
		return err
	}

	db, err := dburl.Open(os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer db.Close()

	row := db.QueryRow("SELECT group_id FROM users WHERE id = $1", ins.Token.UserId)
	err = row.Scan(&ins.GroupId)
	if err != nil {
		return err
	}

	return nil
}
