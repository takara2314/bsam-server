package racing

import (
	"bsam-server/v3/auth"
	"log"
)

// auth authorizes the client.
func (c *Client) auth(msg *AuthInfo) {
	if ok := auth.VerifyJWT(msg.Token); !ok {
		log.Println("Unauthorized:", c.ID)
		c.sendFailedAuthMsg()
		c.Hub.Unregister <- c
		return
	}

	if msg.Role == "mark" && msg.MarkNo == 0 {
		log.Println("Not selecting mark no:", c.ID)
		c.Hub.Unregister <- c
		return
	}

	// If the client has not linked yet, link it.
	if oldID := c.Hub.findClientID(c.UserID); oldID == "" {
		c.link(msg.UserID, msg.Role, msg.MarkNo)
	} else {
		c.restore(oldID)
	}

	c.sendFirstAnnounce()
}

// link links the client.
func (c *Client) link(userID string, role string, markNo int) {
	c.UserID = userID
	c.Role = role

	switch c.Role {
	case "athlete":
		c.Hub.Athletes[c.ID] = c
	case "mark":
		c.MarkNo = markNo
		c.Hub.Marks[c.ID] = c
	case "manager":
		c.Hub.Managers[c.ID] = c
	}

	log.Printf("Linked: %s <=> %s (%s)\n", c.ID, c.UserID, c.Role)

	// Send the authorize result message
	c.sendNewAuthMsg()
}

// restore restores the client.
func (c *Client) restore(oldID string) {
	oldClient := c.Hub.Clients[oldID]

	// Switch data from old to new
	c.UserID = oldClient.UserID
	c.Role = oldClient.Role
	c.MarkNo = oldClient.MarkNo
	c.NextMarkNo = oldClient.NextMarkNo
	c.CourseLimit = oldClient.CourseLimit
	c.Location = oldClient.Location

	// Delete the old client instance
	delete(c.Hub.Clients, oldID)
	delete(c.Hub.Athletes, oldID)
	delete(c.Hub.Marks, oldID)
	delete(c.Hub.Managers, oldID)

	log.Printf("Restored: %s <=> %s (%s)\n", c.ID, c.UserID, c.Role)

	// Send the authorize result message
	c.sendRestoreAuthMsg()
}

// sendFailedAuthMsg sends the failed authorize result message.
func (c *Client) sendFailedAuthMsg() {
	c.sendAuthResultMsgEvent(&AuthResultMsg{
		Authed:   false,
		LinkType: "failed",
	})
}

// sendNewAuthMsg sends the authorize result message for the newbie.
func (c *Client) sendNewAuthMsg() {
	c.sendAuthResultMsgEvent(&AuthResultMsg{
		Authed:   true,
		UserID:   c.UserID,
		Role:     c.Role,
		MarkNo:   c.MarkNo,
		LinkType: "new",
	})
}

// sendRestoreAuthMsg sends the authorize result message for the restored.
func (c *Client) sendRestoreAuthMsg() {
	c.sendAuthResultMsgEvent(&AuthResultMsg{
		Authed:   true,
		UserID:   c.UserID,
		Role:     c.Role,
		MarkNo:   c.MarkNo,
		LinkType: "restore",
	})
}

// sendFirstAnnounce sends the first announce message.
func (c *Client) sendFirstAnnounce() {
	switch c.Role {
	case "athlete":
		c.sendMarkPosMsg()
		c.sendStartRaceMsg()
	case "manager":
		c.sendLiveMsg()
		c.sendStartRaceMsg()
	}
}
