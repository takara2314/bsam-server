//nolint:testpackage // These tests validate internal session-dedup behavior directly.
package racing

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	testAthleteUserID = "athlete1"
	testMarkUserID    = "mark1"
	testJWTSecret     = "test-secret"
)

func TestMain(m *testing.M) {
	_ = os.Setenv("JWT_SECRET", testJWTSecret)

	os.Exit(m.Run())
}

func TestAuthReplacesActiveAthleteSession(t *testing.T) {
	t.Parallel()

	hub := newTestHub()
	token := newTestToken(t)

	oldClient := newTestClient(hub, "athlete-old")
	oldClient.UserID = testAthleteUserID
	oldClient.Role = AthleteRole
	oldClient.NextMarkNo = 3
	oldClient.CourseLimit = 12.5
	oldClient.BatteryLevel = 66
	oldClient.Location = Location{Lat: 34.0, Lng: 136.0, Acc: 2.5}
	oldClient.ConnectedAt = time.Now().Add(-time.Minute)
	hub.Clients[oldClient.ID] = oldClient
	hub.Athletes[oldClient.ID] = oldClient

	newClient := newTestClient(hub, "athlete-new")
	hub.Clients[newClient.ID] = newClient

	newClient.auth(&AuthInfo{
		Token:  token,
		UserID: testAthleteUserID,
		Role:   AthleteRole,
	})
	drainHubLifecycle(hub)

	if hub.Athletes[newClient.ID] == nil {
		t.Fatalf("expected new athlete session to be registered")
	}

	if _, exists := hub.Clients[oldClient.ID]; exists {
		t.Fatalf("expected old athlete session to be removed")
	}

	if newClient.NextMarkNo != 3 {
		t.Fatalf("expected next mark number to be restored, got %d", newClient.NextMarkNo)
	}
}

func TestAuthReplacesActiveMarkSession(t *testing.T) {
	t.Parallel()

	hub := newTestHub()
	token := newTestToken(t)

	oldClient := newTestClient(hub, "mark-old")
	oldClient.UserID = testMarkUserID
	oldClient.Role = MarkRole
	oldClient.MarkNo = 1
	oldClient.BatteryLevel = 90
	oldClient.Location = Location{
		Lat:            34.0,
		Lng:            136.0,
		Acc:            0.0,
		PositionSource: PositionSourceManual,
	}
	oldClient.ConnectedAt = time.Now().Add(-time.Minute)
	hub.Clients[oldClient.ID] = oldClient
	hub.Marks[oldClient.ID] = oldClient

	newClient := newTestClient(hub, "mark-new")
	hub.Clients[newClient.ID] = newClient

	newClient.auth(&AuthInfo{
		Token:  token,
		UserID: testMarkUserID,
		Role:   MarkRole,
		MarkNo: 1,
	})
	drainHubLifecycle(hub)

	if hub.Marks[newClient.ID] == nil {
		t.Fatalf("expected new mark session to be registered")
	}

	if _, exists := hub.Clients[oldClient.ID]; exists {
		t.Fatalf("expected old mark session to be removed")
	}

	if newClient.MarkNo != 1 {
		t.Fatalf("expected mark number to be restored, got %d", newClient.MarkNo)
	}

	if newClient.Location.PositionSource != PositionSourceManual {
		t.Fatalf("expected position source to be carried over, got %q", newClient.Location.PositionSource)
	}
}

func TestGenerateLiveMsgChoosesNewestDuplicateSessions(t *testing.T) {
	t.Parallel()

	hub := newTestHub()
	older := time.Now().Add(-time.Minute)
	newer := time.Now()

	oldAthlete := newTestClient(hub, "athlete-old")
	oldAthlete.UserID = testAthleteUserID
	oldAthlete.Role = AthleteRole
	oldAthlete.NextMarkNo = 1
	oldAthlete.ConnectedAt = older
	hub.Athletes[oldAthlete.ID] = oldAthlete

	newAthlete := newTestClient(hub, "athlete-new")
	newAthlete.UserID = testAthleteUserID
	newAthlete.Role = AthleteRole
	newAthlete.NextMarkNo = 2
	newAthlete.ConnectedAt = newer
	hub.Athletes[newAthlete.ID] = newAthlete

	oldMark := newTestClient(hub, "mark-old")
	oldMark.UserID = testMarkUserID
	oldMark.Role = MarkRole
	oldMark.MarkNo = 1
	oldMark.ConnectedAt = older
	oldMark.Location = Location{Acc: 2.0, PositionSource: PositionSourceGPS}
	hub.Marks[oldMark.ID] = oldMark

	newMark := newTestClient(hub, "mark-new")
	newMark.UserID = testMarkUserID
	newMark.Role = MarkRole
	newMark.MarkNo = 1
	newMark.ConnectedAt = newer
	newMark.Location = Location{Acc: 0.0, PositionSource: PositionSourceManual}
	hub.Marks[newMark.ID] = newMark

	msg := hub.generateLiveMsg()

	if len(msg.Athletes) != 1 {
		t.Fatalf("expected one deduplicated athlete, got %d", len(msg.Athletes))
	}

	if msg.Athletes[0].NextMarkNo != 2 {
		t.Fatalf("expected newest athlete state, got %d", msg.Athletes[0].NextMarkNo)
	}

	if msg.Marks[0].Position.PositionSource != PositionSourceManual {
		t.Fatalf("expected newest mark state, got %q", msg.Marks[0].Position.PositionSource)
	}
}

func TestInsertTypeToJSONIncludesPositionSource(t *testing.T) {
	t.Parallel()

	payload := insertTypeToJSON(&MarkPosMsg{
		MarkNum: 1,
		Marks: []Mark{
			{
				MarkNo: 1,
				Position: Position{
					Lat:            34.0,
					Lng:            136.0,
					Acc:            0.0,
					PositionSource: PositionSourceManual,
				},
			},
		},
	}, "mark_position")

	var decoded map[string]any

	err := json.Unmarshal(payload, &decoded)
	if err != nil {
		t.Fatalf("failed to decode payload: %v", err)
	}

	marks, ok := decoded["marks"].([]any)
	if !ok {
		t.Fatalf("expected marks to be []any, got %T", decoded["marks"])
	}

	firstMark, ok := marks[0].(map[string]any)
	if !ok {
		t.Fatalf("expected first mark to be map[string]any, got %T", marks[0])
	}

	position, ok := firstMark["position"].(map[string]any)
	if !ok {
		t.Fatalf("expected position to be map[string]any, got %T", firstMark["position"])
	}

	if position["position_source"] != PositionSourceManual {
		t.Fatalf("expected position_source to be serialized, got %#v", position["position_source"])
	}
}

func newTestHub() *Hub {
	return &Hub{
		AssociationID: "assoc",
		Clients:       make(map[string]*Client),
		Athletes:      make(map[string]*Client),
		Marks:         make(map[string]*Client),
		Managers:      make(map[string]*Client),
		Disconnectors: make(map[string]*Client),
		MarkNum:       MarkNum,
		Register:      make(chan *Client, 10),
		Disconnect:    make(chan *Client, 10),
		Unregister:    make(chan *Client, 10),
	}
}

func newTestClient(hub *Hub, id string) *Client {
	return &Client{
		ID:           id,
		ConnectedAt:  time.Now(),
		Hub:          hub,
		MarkNo:       -1,
		NextMarkNo:   1,
		Location:     Location{},
		BatteryLevel: -1,
		Send:         make(chan []byte, 10),
	}
}

func drainHubLifecycle(hub *Hub) {
	for {
		select {
		case client := <-hub.Unregister:
			hub.unregisterEvent(client)
		case client := <-hub.Disconnect:
			hub.disconnectEvent(client)
		default:
			return
		}
	}
}

func newTestToken(t *testing.T) string {
	t.Helper()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "tester",
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	signedToken, err := token.SignedString([]byte(testJWTSecret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	return signedToken
}
