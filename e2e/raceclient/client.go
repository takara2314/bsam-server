package raceclient

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// 受信メッセージ (ingress) の最大サイズ: 1KB
	maxIngressMessageBytes = 1024
)

type Client struct {
	DeviceID  string
	Conn      *websocket.Conn
	Resp      *http.Response
	url       url.URL
	receiveCh chan []byte
	closeCh   chan struct{}
	closeOnce sync.Once
	mu        sync.RWMutex
}

var (
	ErrClientClosed = errors.New("クライアントは閉じられています")
)

func NewClient(u url.URL, deviceID string) *Client {
	return &Client{
		DeviceID:  deviceID,
		url:       u,
		receiveCh: make(chan []byte, maxIngressMessageBytes),
		closeCh:   make(chan struct{}),
	}
}

func (c *Client) Connect(ctx context.Context, timeout time.Duration) error {
	dialer := websocket.Dialer{
		Proxy:            websocket.DefaultDialer.Proxy,
		HandshakeTimeout: timeout,
	}

	conn, resp, err := dialer.DialContext(ctx, c.url.String(), nil)
	if err != nil {
		return fmt.Errorf("このURLに接続できませんでした: %w", err)
	}

	c.mu.Lock()
	c.Conn = conn
	c.Resp = resp
	c.mu.Unlock()

	go c.readPump()

	return nil
}

func (c *Client) Close() {
	c.closeOnce.Do(func() {
		close(c.closeCh)
		c.mu.Lock()
		if c.Conn != nil {
			c.Conn.Close()
		}
		c.mu.Unlock()
	})
}
