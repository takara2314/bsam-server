package auth

import (
	"log"
	"net/http"
	"os"
	"sailing-assist-mie-api/abort"
	"sailing-assist-mie-api/bsamdb"
	"sailing-assist-mie-api/inspector"
	"sailing-assist-mie-api/message"
	"sailing-assist-mie-api/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type PasswordPOSTJSON struct {
	LoginID  string `json:"login_id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}

// passwordPOST is /auth/password POST request handler.
func passwordPOST(c *gin.Context) {
	ins := inspector.Inspector{Request: c.Request}

	// Only JSON.
	if !ins.IsJSON() {
		abort.BadRequest(c, message.OnlyJSON)
		return
	}

	var json PasswordPOSTJSON

	// Check all of the require field is not blanked.
	err := c.ShouldBindJSON(&json)
	if err != nil {
		abort.BadRequest(c, message.NotMeetAllRequest)
		return
	}

	// Connect to the database
	db, err := bsamdb.Open()
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer db.DB.Close()

	row := db.DB.QueryRow("SELECT id FROM users WHERE login_id = $1 AND password = digest($2, 'sha3-256')", json.LoginID, json.Password)
	var userID string
	err = row.Scan(&userID)
	if err != nil {
		abort.Forbidden(c, message.WrongIDorPassword)
		return
	}

	// Generate JWT token.
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 30 * 3).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Println(err)
		panic(err)
	}

	_, err = db.Insert(
		"tokens",
		[]bsamdb.Field{
			{
				Column: "token",
				Value:  tokenString,
			},
			{
				Column:  "permissions",
				Value2d: utils.StrSliceToAnySlice([]string{"*"}),
			},
			{
				Column: "user_id",
				Value:  userID,
			},
			{
				Column: "description",
				Value:  "user token",
			},
		},
	)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	c.JSON(http.StatusOK, TokenResponse{
		Status: "OK",
		Token:  tokenString,
	})
}
