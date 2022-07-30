package auth

import (
	"bsam-server/v1/abort"
	"bsam-server/v1/bsamdb"
	"bsam-server/v1/inspector"
	"bsam-server/v1/message"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type TokenPOSTJSON struct {
	Token string `json:"token" binding:"required"`
}

// tokenPOST is /auth/token POST request handler.
func tokenPOST(c *gin.Context) {
	ins := inspector.Inspector{Request: c.Request}

	// Only JSON.
	if !ins.IsJSON() {
		abort.BadRequest(c, message.OnlyJSON)
		return
	}

	var json TokenPOSTJSON

	// Check all of the require field is not blanked.
	err := c.ShouldBindJSON(&json)
	if err != nil {
		abort.BadRequest(c, message.NotMeetAllRequest)
		return
	}

	token, err := jwt.Parse(json.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if token == nil {
		abort.BadRequest(c, message.InformedJWT)
		return
	}

	if !token.Valid {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				abort.BadRequest(c, message.InformedJWT)
				return
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				abort.Forbidden(c, message.ExpiredOrNotValidYetJWT)
				return
			}
		}
		abort.Forbidden(c, message.InvalidJWT)
	}

	// Connect to the database
	db, err := bsamdb.Open()
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer db.DB.Close()

	ok, err := db.IsExist("tokens", "token", json.Token)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	if !ok {
		abort.Forbidden(c, message.RemovedJWT)
		return
	}

	abort.OK(c, message.ValidJWT)
	return
}