package racehub

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/oklog/ulid/v2"
)

const (
	writeWaitSec       = 10 * time.Second
	pongWaitSec        = 60 * time.Second
	pingPeriodSec      = (pongWaitSec * 9) / 10
	maxMessageSizeByte = 512
)

type Client struct {
	ID   string
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
}

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (c *Client) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("id", c.ID),
		slog.String("address", c.Conn.RemoteAddr().String()),
		slog.String("assoc_id", c.Hub.AssocID),
	)
}

func (h *Hub) Register(conn *websocket.Conn) *Client {
	id := ulid.Make().String()

	client := &Client{
		ID:   id,
		Hub:  h,
		Conn: conn,
		Send: make(chan []byte, maxMessageSizeByte),
	}

	h.mu.Lock()
	h.clients[id] = client
	h.mu.Unlock()

	go client.readPump()
	go client.writePump()

	return client
}

func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exist := h.clients[c.ID]; exist {
		delete(h.clients, c.ID)
		close(c.Send)
	}
}
