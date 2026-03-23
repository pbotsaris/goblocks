package socketmode

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/coder/websocket"
)

// writeRequest represents a request to write data to the WebSocket.
type writeRequest struct {
	data   []byte
	result chan error
}

// writer provides thread-safe writes to a WebSocket connection.
// All writes are serialized through a single goroutine to avoid
// concurrent write issues with the WebSocket library.
type writer struct {
	conn         *websocket.Conn
	queue        chan writeRequest
	writeTimeout time.Duration
	metrics      MetricsHook

	mu        sync.Mutex
	closed    bool
	closeOnce sync.Once
	done      chan struct{}
}

// writerConfig configures the writer behavior.
type writerConfig struct {
	queueSize    int
	writeTimeout time.Duration
	metrics      MetricsHook
}

// defaultWriterConfig returns sensible defaults.
func defaultWriterConfig() writerConfig {
	return writerConfig{
		queueSize:    100,
		writeTimeout: 3 * time.Second,
		metrics:      &NoopMetrics{},
	}
}

// newWriter creates a new writer for the given WebSocket connection.
func newWriter(conn *websocket.Conn, cfg writerConfig) *writer {
	if cfg.queueSize <= 0 {
		cfg.queueSize = 100
	}
	if cfg.writeTimeout <= 0 {
		cfg.writeTimeout = 3 * time.Second
	}
	if cfg.metrics == nil {
		cfg.metrics = &NoopMetrics{}
	}

	w := &writer{
		conn:         conn,
		queue:        make(chan writeRequest, cfg.queueSize),
		writeTimeout: cfg.writeTimeout,
		metrics:      cfg.metrics,
		done:         make(chan struct{}),
	}

	go w.loop()
	return w
}

// loop is the main write loop that drains the queue and writes to the socket.
func (w *writer) loop() {
	defer close(w.done)

	for req := range w.queue {
		ctx, cancel := context.WithTimeout(context.Background(), w.writeTimeout)
		err := w.conn.Write(ctx, websocket.MessageText, req.data)
		cancel()

		if req.result != nil {
			req.result <- err
			close(req.result)
		}
	}
}

// Write sends data to the WebSocket, blocking until the write completes or times out.
// Returns an error if the write fails or the context is cancelled.
func (w *writer) Write(ctx context.Context, data []byte) error {
	// Check if closed before attempting to send
	w.mu.Lock()
	if w.closed {
		w.mu.Unlock()
		return ErrConnectionClosed
	}
	w.mu.Unlock()

	result := make(chan error, 1)
	req := writeRequest{
		data:   data,
		result: result,
	}

	// Report queue depth
	w.metrics.WriteQueueDepth(len(w.queue))

	select {
	case w.queue <- req:
		// Request queued, wait for result
	case <-ctx.Done():
		return ctx.Err()
	case <-w.done:
		return ErrConnectionClosed
	}

	select {
	case err := <-result:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case <-w.done:
		return ErrConnectionClosed
	}
}

// WriteJSON marshals v to JSON and sends it to the WebSocket.
func (w *writer) WriteJSON(ctx context.Context, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return w.Write(ctx, data)
}

// WriteAck sends an acknowledgment for the given envelope.
func (w *writer) WriteAck(ctx context.Context, envelopeID string, payload any) error {
	ack := Ack{
		EnvelopeID: envelopeID,
		Payload:    payload,
	}
	return w.WriteJSON(ctx, ack)
}

// Close stops the writer and waits for pending writes to complete.
func (w *writer) Close() {
	w.closeOnce.Do(func() {
		w.mu.Lock()
		w.closed = true
		w.mu.Unlock()
		close(w.queue)
	})
	<-w.done
}

// Ping sends a WebSocket ping frame.
func (w *writer) Ping(ctx context.Context) error {
	return w.conn.Ping(ctx)
}
