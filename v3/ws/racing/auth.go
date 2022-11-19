package racing

import "log"

func (c *Client) auth(msg *AuthInfo) {
	userID, role, err := getUserInfoFromJWT(msg.Token)
	if err != nil {
		log.Println("Unauthorized:", c.ID)
		c.Hub.Unregister <- c
		return
	}

	if role == "mark" && msg.MarkNo == 0 {
		log.Println("Not select mark no:", c.ID)
		c.Hub.Unregister <- c
		return
	}

	log.Printf("Linked: %s <=> %s (%s)\n", c.ID, userID, role)

	c.UserID = userID
	c.Role = role

	switch role {
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
