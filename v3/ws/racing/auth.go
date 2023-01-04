package racing

import (
	"bsam-server/v3/auth"
	"log"
)

// auth authorizes the client.
func (c *Client) auth(msg *AuthInfo) {
	if ok := auth.VerifyJWT(msg.Token); !ok {
		log.Println("Unauthorized:", c.ID)
		c.Hub.Unregister <- c
		return
	}

	if msg.Role == "mark" && msg.MarkNo == 0 {
		log.Println("Not selecting mark no:", c.ID)
		c.Hub.Unregister <- c
		return
	}

	c.UserID = msg.UserID
	c.Role = msg.Role

	log.Printf("Linked: %s <=> %s (%s)\n", c.ID, c.UserID, c.Role)

	switch c.Role {
	case "athlete":
		c.Hub.Athletes[c.ID] = c
		c.sendMarkPosMsg()
		c.sendStartRaceMsg()
	case "mark":
		c.MarkNo = msg.MarkNo
		c.Hub.Marks[c.ID] = c
	case "manage":
		c.sendLiveMsg()
		c.sendStartRaceMsg()
	}
}
