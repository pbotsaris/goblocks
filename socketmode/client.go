package socketmode

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

// EventHandler processes an envelope and returns a type-safe response.
type EventHandler func(ctx context.Context, envelope *Envelope) Response

// Client is a Socket Mode client for Slack.
type Client struct {
	appToken   string
	httpClient *http.Client
	logger     *slog.Logger
	metrics    MetricsHook
	backoff    *Backoff

	// Handler configuration
	handlers       map[string]EventHandler
	handlersMu     sync.RWMutex
	maxConcurrency int
	handlerTimeout time.Duration

	// Connection configuration
	helloTimeout time.Duration
	pingInterval time.Duration
	pongTimeout  time.Duration

	// Shutdown coordination
	shutdownMu sync.Mutex
	shutting   bool
}

// Option configures the Client.
type Option func(*Client)

// New creates a new Socket Mode client.
func New(appToken string, opts ...Option) *Client {
	c := &Client{
		appToken:       appToken,
		httpClient:     &http.Client{Timeout: 30 * time.Second},
		logger:         slog.Default(),
		metrics:        &NoopMetrics{},
		backoff:        NewBackoff(DefaultBackoffConfig()),
		handlers:       make(map[string]EventHandler),
		maxConcurrency: 10,
		handlerTimeout: 30 * time.Second,
		helloTimeout:   30 * time.Second,
		pingInterval:   5 * time.Second,
		pongTimeout:    10 * time.Second,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// WithLogger sets the logger for the client.
func WithLogger(logger *slog.Logger) Option {
	return func(c *Client) {
		if logger != nil {
			c.logger = logger
		}
	}
}

// WithMetrics sets the metrics hook for the client.
func WithMetrics(hook MetricsHook) Option {
	return func(c *Client) {
		if hook != nil {
			c.metrics = hook
		}
	}
}

// WithHTTPClient sets the HTTP client for API calls.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		if client != nil {
			c.httpClient = client
		}
	}
}

// WithMaxConcurrency sets the maximum number of concurrent handlers.
func WithMaxConcurrency(n int) Option {
	return func(c *Client) {
		if n > 0 {
			c.maxConcurrency = n
		}
	}
}

// WithHandlerTimeout sets the timeout for handler execution.
func WithHandlerTimeout(d time.Duration) Option {
	return func(c *Client) {
		if d > 0 {
			c.handlerTimeout = d
		}
	}
}

// WithHelloTimeout sets the timeout for waiting for the hello message.
func WithHelloTimeout(d time.Duration) Option {
	return func(c *Client) {
		if d > 0 {
			c.helloTimeout = d
		}
	}
}

// WithPingInterval sets the interval between ping messages.
func WithPingInterval(d time.Duration) Option {
	return func(c *Client) {
		if d > 0 {
			c.pingInterval = d
		}
	}
}

// WithPongTimeout sets the timeout for receiving pong responses.
func WithPongTimeout(d time.Duration) Option {
	return func(c *Client) {
		if d > 0 {
			c.pongTimeout = d
		}
	}
}

// On registers a handler for the given event type.
func (c *Client) On(eventType string, handler EventHandler) {
	c.handlersMu.Lock()
	defer c.handlersMu.Unlock()
	c.handlers[eventType] = handler
}

// OnSlashCommand registers a handler for slash commands.
func (c *Client) OnSlashCommand(handler EventHandler) {
	c.On(EnvelopeTypeSlashCommands, handler)
}

// OnInteractive registers a handler for interactive events (buttons, modals, etc).
func (c *Client) OnInteractive(handler EventHandler) {
	c.On(EnvelopeTypeInteractive, handler)
}

// OnEventsAPI registers a handler for Events API events.
func (c *Client) OnEventsAPI(handler EventHandler) {
	c.On(EnvelopeTypeEventsAPI, handler)
}

// Run starts the client and maintains a connection to Slack.
// It reconnects automatically on disconnection until ctx is cancelled.
// Returns a permanent error if the connection cannot be established.
func (c *Client) Run(ctx context.Context) error {
	for {
		err := c.runOnce(ctx)

		// Check if we should stop
		if ctx.Err() != nil {
			return ctx.Err()
		}

		// Check for permanent errors
		if IsPermanentError(err) {
			c.logger.Error("permanent error, not retrying", "error", err)
			return err
		}

		// Calculate backoff delay
		delay := c.backoff.NextDelay()
		c.logger.Info("reconnecting",
			"error", err,
			"attempt", c.backoff.Attempts(),
			"delay", delay,
		)
		c.metrics.ReconnectAttempt(c.backoff.Attempts(), delay)

		// Wait before reconnecting
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}
}

// runOnce runs a single connection lifecycle.
func (c *Client) runOnce(ctx context.Context) error {
	// Get WebSocket URL
	wsURL, err := c.openConnection(ctx)
	if err != nil {
		return err
	}

	// Generate connection ID
	connID := uuid.New().String()[:8]

	// Establish connection
	conn, err := dial(ctx, connectionConfig{
		id:             connID,
		url:            wsURL,
		helloTimeout:   c.helloTimeout,
		pingInterval:   c.pingInterval,
		pongTimeout:    c.pongTimeout,
		writeQueueSize: 100,
		writeTimeout:   3 * time.Second,
		logger:         c.logger,
		metrics:        c.metrics,
	})
	if err != nil {
		return err
	}
	defer conn.Close()

	// Mark connection as stable for backoff
	c.backoff.MarkConnected()

	// Run the read loop with concurrency control
	return c.readLoop(ctx, conn)
}

// readLoop reads messages from the connection and dispatches to handlers.
func (c *Client) readLoop(ctx context.Context, conn *connection) error {
	// Semaphore for concurrency control
	sem := make(chan struct{}, c.maxConcurrency)

	// Track in-flight handlers for graceful shutdown
	var wg sync.WaitGroup
	defer wg.Wait()

	for {
		// Check for context cancellation
		select {
		case <-ctx.Done():
			c.logger.Debug("context cancelled, shutting down")
			return ctx.Err()
		case <-conn.Done():
			return ErrConnectionClosed
		default:
		}

		// Read next message
		data, err := conn.Read(ctx)
		if err != nil {
			return err
		}

		// Parse envelope type
		var msg struct {
			Type       string `json:"type"`
			EnvelopeID string `json:"envelope_id"`
		}
		if err := json.Unmarshal(data, &msg); err != nil {
			c.logger.Warn("invalid json", "error", err)
			continue
		}

		// Handle disconnect
		if msg.Type == EnvelopeTypeDisconnect {
			var disconnect DisconnectMessage
			if err := json.Unmarshal(data, &disconnect); err == nil {
				c.logger.Info("disconnect requested",
					"reason", disconnect.Reason,
				)
				if disconnect.Reason == DisconnectReasonLinkDisabled {
					return &PermanentError{Message: "socket mode disabled"}
				}
			}
			return &RetryableError{Message: "disconnect requested"}
		}

		// Skip hello (already handled in dial)
		if msg.Type == EnvelopeTypeHello {
			continue
		}

		// Parse full envelope
		var env Envelope
		if err := json.Unmarshal(data, &env); err != nil {
			c.logger.Warn("invalid envelope", "error", err)
			continue
		}

		c.metrics.EnvelopeReceived(env.Type)
		receivedAt := time.Now()

		// Get handler
		c.handlersMu.RLock()
		handler, ok := c.handlers[env.Type]
		c.handlersMu.RUnlock()

		if !ok {
			// No handler, just ack
			c.ackEnvelope(ctx, conn, &env, nil, receivedAt)
			continue
		}

		// Handle based on whether response payload is accepted
		if env.AcceptsResponsePayload {
			// Must run handler synchronously to include response in ack
			wg.Add(1)
			go func(env Envelope) {
				defer wg.Done()

				// Acquire semaphore
				select {
				case sem <- struct{}{}:
					defer func() { <-sem }()
				case <-ctx.Done():
					return
				}

				resp := c.safeHandler(ctx, &env, handler)
				c.ackEnvelope(ctx, conn, &env, resp, receivedAt)
			}(env)
		} else {
			// Ack immediately, run handler async
			c.ackEnvelope(ctx, conn, &env, nil, receivedAt)

			wg.Add(1)
			go func(env Envelope) {
				defer wg.Done()

				// Acquire semaphore
				select {
				case sem <- struct{}{}:
					defer func() { <-sem }()
				case <-ctx.Done():
					return
				}

				c.safeHandler(ctx, &env, handler)
			}(env)
		}

		// Periodically check if backoff should be reset
		c.backoff.CheckStable()
	}
}

// safeHandler executes a handler with panic recovery and timeout.
func (c *Client) safeHandler(ctx context.Context, env *Envelope, handler EventHandler) (resp Response) {
	c.metrics.HandlerStarted(env.Type)
	start := time.Now()

	defer func() {
		duration := time.Since(start)

		if r := recover(); r != nil {
			c.logger.Error("handler panic",
				"panic", r,
				"envelope_id", env.EnvelopeID,
				"type", env.Type,
			)
			c.metrics.HandlerPanic(env.Type, r)
			resp = EmptyResponse{}
		}

		c.metrics.HandlerCompleted(env.Type, duration, nil)
	}()

	// Create context with timeout
	handlerCtx, cancel := context.WithTimeout(ctx, c.handlerTimeout)
	defer cancel()

	return handler(handlerCtx, env)
}

// ackEnvelope sends an acknowledgment for the envelope.
func (c *Client) ackEnvelope(ctx context.Context, conn *connection, env *Envelope, resp Response, receivedAt time.Time) {
	var payload any
	if resp != nil {
		payload = resp.toPayload()
	}

	ackCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := conn.Writer().WriteAck(ackCtx, env.EnvelopeID, payload); err != nil {
		c.logger.Error("failed to ack",
			"envelope_id", env.EnvelopeID,
			"error", err,
		)
		return
	}

	latency := time.Since(receivedAt)
	c.metrics.EnvelopeAcked(env.Type, latency)
	c.logger.Debug("acked envelope",
		"envelope_id", env.EnvelopeID,
		"type", env.Type,
		"latency", latency,
	)
}

// openConnection calls apps.connections.open to get a WebSocket URL.
func (c *Client) openConnection(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "POST",
		"https://slack.com/api/apps.connections.open", nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.appToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", ClassifyNetworkError(err)
	}
	defer resp.Body.Close()

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.logger.Debug("apps.connections.open failed",
			"status", resp.StatusCode,
			"body", string(body),
		)
		return "", ClassifyHTTPError(resp.StatusCode)
	}

	var result ConnectionOpenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decoding response: %w", err)
	}

	if !result.OK {
		return "", ClassifySlackError(result.Error)
	}

	c.logger.Debug("got websocket url", "url", result.URL)
	return result.URL, nil
}
