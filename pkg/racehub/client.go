package racehub

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/oklog/ulid/v2"
)

const (
	// WebSocketの書き込み操作のタイムアウト時間: 1秒
	// 書き込みが完了しない場合はエラーとする
	writeWaitSec = 1 * time.Second

	// クライアントからのPong応答を待つ最大時間: 10秒
	// Pongが返ってこない場合は接続に問題があると判断する
	pongWaitSec = 10 * time.Second

	// クライアントへのマークの位置情報送信間隔: 5秒
	sendingMarkGeolocationsTickerPeriodSec = 5 * time.Second

	// サーバーからクライアントへPingを送信する間隔: 9秒
	// タイムアウト前に必ずPingを送信する
	pingPeriodSec = (pongWaitSec * 9) / 10

	// 受信メッセージ (ingress) の最大サイズ: 1KB
	// これより大きいメッセージは拒否される
	maxIngressMessageBytes = 1024
)

type Client struct {
	ID                string
	Hub               *Hub
	Conn              *websocket.Conn
	Send              chan []byte
	StoppingWritePump chan bool

	DeviceID   string
	Role       string
	MarkNo     int
	NextMarkNo int
}

// WebSocketアップグレーダー: HTTP接続をWebSocket接続にアップグレードする設定
var Upgrader = websocket.Upgrader{
	// 受信用バッファサイズ: 4KB
	// 大きめのバッファでネットワーク操作の効率を高める
	ReadBufferSize: 4096,

	// 送信用バッファサイズ: 4KB
	// 大きめのバッファでネットワーク操作の効率を高める
	WriteBufferSize: 4096,

	// WebSocketの圧縮を有効化
	// データ転送量を削減するが、CPUの使用量が若干増加する可能性がある
	EnableCompression: true,

	// クロスオリジンリクエストのチェック
	// 本番環境では適切なオリジン検証を実装すべき
	// 現在の設定は全てのオリジンを許可 (開発用)
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (c *Client) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("id", c.ID),
		slog.String("address", c.Conn.RemoteAddr().String()),
		slog.String("assoc_id", c.Hub.AssocID),
		slog.String("device_id", c.DeviceID),
		slog.String("role", c.Role),
		slog.Int("mark_no", c.MarkNo),
		slog.Int("next_mark_no", c.NextMarkNo),
	)
}

func (h *Hub) Register(conn *websocket.Conn) *Client {
	id := ulid.Make().String()

	client := &Client{
		ID:                id,
		Hub:               h,
		Conn:              conn,
		Send:              make(chan []byte, maxIngressMessageBytes),
		StoppingWritePump: make(chan bool),

		DeviceID:   "unknown",
		Role:       "unknown",
		MarkNo:     -1,
		NextMarkNo: -1,
	}

	h.Mu.Lock()
	h.clients[id] = client
	h.Mu.Unlock()

	slog.Info(
		"client registered",
		"client", client,
	)

	client.SetPingPongHandler()
	go client.readPump()
	go client.writePump()

	return client
}

func (h *Hub) Unregister(c *Client) {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	slog.Info(
		"client unregistered",
		"client", c,
	)

	if _, exist := h.clients[c.ID]; exist {
		c.StoppingWritePump <- true

		delete(h.clients, c.ID)
		close(c.Send)
	}
}

func (c *Client) SetPingPongHandler() {
	// クライアントからのPingメッセージを処理するハンドラ
	c.Conn.SetPingHandler(func(data string) error {
		slog.Info(
			"WebSocket ping received from client",
			"client", c,
			"data", data,
		)

		return c.Conn.WriteControl(
			websocket.PongMessage,
			[]byte(data),
			time.Now().Add(writeWaitSec),
		)
	})

	// クライアントからのPongメッセージを処理するハンドラ
	c.Conn.SetPongHandler(func(data string) error {
		slog.Info(
			"WebSocket pong received from client",
			"client", c,
			"data", data,
		)

		return c.Conn.SetReadDeadline(time.Now().Add(pongWaitSec))
	})
}
