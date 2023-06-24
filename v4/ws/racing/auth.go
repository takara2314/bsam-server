package racing

import (
	"bsam-server/utils"
	"bsam-server/v4/auth"
	"log"

	"golang.org/x/exp/slices"
)

// auth authorizes the client.
func (c *Client) auth(msg *AuthInfo) {
	// If the client is a guest, not need to verify the token
	if msg.Role == "guest" {
		c.link(utils.RandString(8), "guest", 0)
		c.sendFirstAnnounce()
		return
	}

	if ok := auth.VerifyJWT(msg.Token); !ok {
		log.Println("Unauthorized:", c.ID)
		c.sendFailedAuthMsg()
		c.Hub.Unregister <- c
		return
	}

	if ok := c.isValidRole(); !ok {
		log.Println("Invalid role:", c.ID)
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
	if oldID := c.Hub.findDisconnectedID(msg.UserID); oldID == "" {
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

	if c.Role == "mark" {
		c.MarkNo = markNo
	}

	// Register the client to the role group
	c.registerRoleGroup()

	log.Printf("Linked: %s <=> %s (%s)\n", c.ID, c.UserID, c.Role)

	// Send the authorize result message
	c.sendNewAuthMsg()
}

// restore restores the client.
func (c *Client) restore(oldID string) {
	oldClient := c.Hub.Disconnectors[oldID]

	// Switch data from old to new
	c.UserID = oldClient.UserID
	c.Role = oldClient.Role
	c.MarkNo = oldClient.MarkNo
	c.NextMarkNo = oldClient.NextMarkNo
	c.CourseLimit = oldClient.CourseLimit
	c.Location = oldClient.Location

	// Delete the old client instance
	c.Hub.Unregister <- oldClient

	// Register the new client instance to the role group
	c.registerRoleGroup()

	log.Printf("Restored: %s <=> %s (%s)\n", c.ID, c.UserID, c.Role)

	// Send the authorize result message
	c.sendRestoreAuthMsg()
}

// registerRoleGroup registers the client to the role group.
func (c *Client) registerRoleGroup() {
	switch c.Role {
	case "athlete":
		c.Hub.Athletes[c.ID] = c
	case "mark":
		c.Hub.Marks[c.ID] = c
	case "manager":
		c.Hub.Managers[c.ID] = c
	}
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
		Authed:     true,
		UserID:     c.UserID,
		Role:       c.Role,
		MarkNo:     c.MarkNo,
		NextMarkNo: c.NextMarkNo,
		LinkType:   "restore",
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
	case "guest":
		c.sendLiveMsg()
		c.sendStartRaceMsg()
	}
}

func (c *Client) isValidRole() bool {
	if ok := slices.Contains([]string{"athlete", "mark", "manager"}, c.Role); !ok {
		return false
	}

	return true
}
