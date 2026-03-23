package socketmode

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/coder/websocket"
)

// mustJSON marshals v to JSON, failing the test on error.
func mustJSON(t *testing.T, v any) []byte {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("mustJSON: %v", err)
	}
	return data
}

// mustUnmarshal unmarshals data into v, failing the test on error.
func mustUnmarshal(t *testing.T, data []byte, v any) {
	t.Helper()
	if err := json.Unmarshal(data, v); err != nil {
		t.Fatalf("mustUnmarshal: %v", err)
	}
}

// mockMetrics captures metrics for testing.
type mockMetrics struct {
	mu sync.Mutex

	connectionsOpened int
	connectionsClosed int
	reconnectAttempts int
	envelopesReceived int
	envelopesAcked    int
	handlersStarted   int
	handlersCompleted int
	handlerPanics     int
	pingsSent         int
	pongsReceived     int
	pongTimeouts      int
}

var _ MetricsHook = (*mockMetrics)(nil)

func (m *mockMetrics) ConnectionOpened(connID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.connectionsOpened++
}

func (m *mockMetrics) ConnectionClosed(connID string, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.connectionsClosed++
}

func (m *mockMetrics) ReconnectAttempt(attempt int, delay time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.reconnectAttempts++
}

func (m *mockMetrics) EnvelopeReceived(envType string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.envelopesReceived++
}

func (m *mockMetrics) EnvelopeAcked(envType string, latency time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.envelopesAcked++
}

func (m *mockMetrics) HandlerStarted(envType string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlersStarted++
}

func (m *mockMetrics) HandlerCompleted(envType string, duration time.Duration, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlersCompleted++
}

func (m *mockMetrics) HandlerPanic(envType string, recovered any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlerPanics++
}

func (m *mockMetrics) WriteQueueDepth(depth int) {}

func (m *mockMetrics) PingSent() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pingsSent++
}

func (m *mockMetrics) PongReceived(latency time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pongsReceived++
}

func (m *mockMetrics) PongTimeout() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pongTimeouts++
}

// mockServer is a mock WebSocket server for testing.
type mockServer struct {
	server   *httptest.Server
	handler  func(conn *websocket.Conn)
	upgrader http.Handler
}

// startMockServer starts a mock WebSocket server.
// The handler is called for each connection.
func startMockServer(t *testing.T, handler func(conn *websocket.Conn)) *mockServer {
	t.Helper()

	m := &mockServer{
		handler: handler,
	}

	m.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			t.Logf("websocket accept error: %v", err)
			return
		}
		defer conn.CloseNow()

		if m.handler != nil {
			m.handler(conn)
		}
	}))

	return m
}

// URL returns the WebSocket URL for the server.
func (m *mockServer) URL() string {
	return "ws" + strings.TrimPrefix(m.server.URL, "http")
}

// Close shuts down the server.
func (m *mockServer) Close() {
	m.server.Close()
}

// sendHello sends a hello message on the connection.
func sendHello(t *testing.T, conn *websocket.Conn, appID string, numConnections int) {
	t.Helper()
	hello := HelloMessage{
		Type: "hello",
		ConnectionInfo: ConnectionInfo{
			AppID: appID,
		},
		NumConnections: numConnections,
		DebugInfo: DebugInfo{
			Host:                      "test-host",
			ApproximateConnectionTime: 3600,
		},
	}
	data := mustJSON(t, hello)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := conn.Write(ctx, websocket.MessageText, data); err != nil {
		t.Fatalf("sendHello: %v", err)
	}
}

// sendEnvelope sends an envelope on the connection.
func sendEnvelope(t *testing.T, conn *websocket.Conn, env Envelope) {
	t.Helper()
	data := mustJSON(t, env)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := conn.Write(ctx, websocket.MessageText, data); err != nil {
		t.Fatalf("sendEnvelope: %v", err)
	}
}

// readAck reads and parses an ack from the connection.
func readAck(t *testing.T, conn *websocket.Conn) Ack {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, data, err := conn.Read(ctx)
	if err != nil {
		t.Fatalf("readAck: %v", err)
	}
	var ack Ack
	mustUnmarshal(t, data, &ack)
	return ack
}

// sendDisconnect sends a disconnect message on the connection.
func sendDisconnect(t *testing.T, conn *websocket.Conn, reason string) {
	t.Helper()
	disconnect := DisconnectMessage{
		Type:   "disconnect",
		Reason: reason,
	}
	data := mustJSON(t, disconnect)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := conn.Write(ctx, websocket.MessageText, data); err != nil {
		t.Fatalf("sendDisconnect: %v", err)
	}
}
