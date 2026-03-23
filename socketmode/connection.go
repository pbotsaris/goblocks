package socketmode

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/coder/websocket"
)

// connection represents a single WebSocket connection to Slack.
type connection struct {
	id         string
	url        string
	conn       *websocket.Conn
	writer     *writer
	logger     *slog.Logger
	metrics    MetricsHook
	helloInfo  *HelloMessage
	startTime  time.Time

	// Ping/pong state
	pingInterval  time.Duration
	pongTimeout   time.Duration
	lastPingSent  time.Time
	lastPongRecv  time.Time
	pingMu        sync.Mutex

	// Shutdown coordination
	cancel   context.CancelFunc
	done     chan struct{}
	closeErr error
	closeOnce sync.Once
}

// connectionConfig configures the connection behavior.
type connectionConfig struct {
	id             string
	url            string
	helloTimeout   time.Duration
	pingInterval   time.Duration
	pongTimeout    time.Duration
	writeQueueSize int
	writeTimeout   time.Duration
	logger         *slog.Logger
	metrics        MetricsHook
}

// defaultConnectionConfig returns sensible defaults.
func defaultConnectionConfig() connectionConfig {
	return connectionConfig{
		helloTimeout:   30 * time.Second,
		pingInterval:   5 * time.Second,
		pongTimeout:    10 * time.Second,
		writeQueueSize: 100,
		writeTimeout:   3 * time.Second,
		logger:         slog.Default(),
		metrics:        &NoopMetrics{},
	}
}

// dial establishes a WebSocket connection and waits for the hello message.
func dial(ctx context.Context, cfg connectionConfig) (*connection, error) {
	logger := cfg.logger.With("connection_id", cfg.id)
	logger.Debug("dialing websocket", "url", cfg.url)

	// Establish WebSocket connection
	wsConn, _, err := websocket.Dial(ctx, cfg.url, nil)
	if err != nil {
		return nil, fmt.Errorf("websocket dial: %w", ClassifyNetworkError(err))
	}

	connCtx, cancel := context.WithCancel(ctx)
	c := &connection{
		id:           cfg.id,
		url:          cfg.url,
		conn:         wsConn,
		logger:       logger,
		metrics:      cfg.metrics,
		pingInterval: cfg.pingInterval,
		pongTimeout:  cfg.pongTimeout,
		startTime:    time.Now(),
		cancel:       cancel,
		done:         make(chan struct{}),
	}

	// Create writer
	c.writer = newWriter(wsConn, writerConfig{
		queueSize:    cfg.writeQueueSize,
		writeTimeout: cfg.writeTimeout,
		metrics:      cfg.metrics,
	})

	// Wait for hello message
	helloCtx, helloCancel := context.WithTimeout(connCtx, cfg.helloTimeout)
	defer helloCancel()

	hello, err := c.waitForHello(helloCtx)
	if err != nil {
		c.closeInternal(err)
		return nil, err
	}

	c.helloInfo = hello
	c.lastPongRecv = time.Now() // Initialize pong time

	logger.Info("connection established",
		"app_id", hello.ConnectionInfo.AppID,
		"num_connections", hello.NumConnections,
		"approx_lifetime_seconds", hello.DebugInfo.ApproximateConnectionTime,
	)

	cfg.metrics.ConnectionOpened(cfg.id)

	// Start ping loop
	go c.pingLoop(connCtx)

	return c, nil
}

// waitForHello reads messages until it receives a hello message.
func (c *connection) waitForHello(ctx context.Context) (*HelloMessage, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ErrHelloTimeout
		default:
		}

		_, data, err := c.conn.Read(ctx)
		if err != nil {
			return nil, fmt.Errorf("reading hello: %w", ClassifyNetworkError(err))
		}

		var msg struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(data, &msg); err != nil {
			c.logger.Warn("invalid json while waiting for hello", "error", err)
			continue
		}

		if msg.Type == EnvelopeTypeHello {
			var hello HelloMessage
			if err := json.Unmarshal(data, &hello); err != nil {
				return nil, fmt.Errorf("parsing hello message: %w", err)
			}
			return &hello, nil
		}

		// Ignore other messages before hello (shouldn't happen, but be safe)
		c.logger.Debug("ignoring message before hello", "type", msg.Type)
	}
}

// pingLoop sends periodic pings and monitors for pong responses.
func (c *connection) pingLoop(ctx context.Context) {
	ticker := time.NewTicker(c.pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-c.done:
			return
		case <-ticker.C:
			if err := c.sendPing(ctx); err != nil {
				c.logger.Debug("ping failed", "error", err)
				return
			}
		}
	}
}

// sendPing sends a ping and checks if we've received a recent pong.
func (c *connection) sendPing(ctx context.Context) error {
	c.pingMu.Lock()
	defer c.pingMu.Unlock()

	// Check if last pong is too old
	if !c.lastPongRecv.IsZero() && time.Since(c.lastPongRecv) > c.pongTimeout {
		c.metrics.PongTimeout()
		c.logger.Warn("pong timeout, connection may be dead",
			"last_pong", time.Since(c.lastPongRecv),
			"timeout", c.pongTimeout,
		)
		c.closeInternal(fmt.Errorf("pong timeout"))
		return ErrConnectionClosed
	}

	c.lastPingSent = time.Now()
	c.metrics.PingSent()

	pingCtx, cancel := context.WithTimeout(ctx, c.pingInterval)
	defer cancel()

	if err := c.writer.Ping(pingCtx); err != nil {
		return err
	}

	return nil
}

// recordPong records that a pong was received.
func (c *connection) recordPong() {
	c.pingMu.Lock()
	defer c.pingMu.Unlock()

	latency := time.Since(c.lastPingSent)
	c.lastPongRecv = time.Now()
	c.metrics.PongReceived(latency)
}

// Read reads the next message from the WebSocket.
// Returns the raw message data or an error.
func (c *connection) Read(ctx context.Context) ([]byte, error) {
	_, data, err := c.conn.Read(ctx)
	if err != nil {
		return nil, ClassifyNetworkError(err)
	}
	return data, nil
}

// Writer returns the thread-safe writer for this connection.
func (c *connection) Writer() *writer {
	return c.writer
}

// ID returns the connection ID.
func (c *connection) ID() string {
	return c.id
}

// HelloInfo returns the hello message received on connect.
func (c *connection) HelloInfo() *HelloMessage {
	return c.helloInfo
}

// Done returns a channel that's closed when the connection is closed.
func (c *connection) Done() <-chan struct{} {
	return c.done
}

// Close gracefully closes the connection.
func (c *connection) Close() error {
	return c.closeInternal(nil)
}

// CloseWithError closes the connection with a specific error.
func (c *connection) CloseWithError(err error) error {
	return c.closeInternal(err)
}

func (c *connection) closeInternal(err error) error {
	c.closeOnce.Do(func() {
		c.closeErr = err
		c.cancel()

		// Close writer first (drains queue)
		c.writer.Close()

		// Send close frame
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		closeErr := c.conn.Close(websocket.StatusNormalClosure, "client closing")
		if closeErr != nil {
			c.logger.Debug("error sending close frame", "error", closeErr)
		}

		_ = closeCtx // Used for timeout context

		// Record metrics
		duration := time.Since(c.startTime)
		c.metrics.ConnectionClosed(c.id, duration)
		c.logger.Info("connection closed", "duration", duration, "error", err)

		close(c.done)
	})

	return c.closeErr
}
