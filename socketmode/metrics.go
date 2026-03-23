package socketmode

import "time"

// MetricsHook allows users to collect metrics from the socket mode client.
// Implement this interface and pass it via WithMetrics() to receive callbacks.
// All methods should be non-blocking and safe for concurrent use.
type MetricsHook interface {
	// ConnectionOpened is called when a WebSocket connection is established.
	ConnectionOpened(connID string)

	// ConnectionClosed is called when a WebSocket connection is closed.
	ConnectionClosed(connID string, duration time.Duration)

	// ReconnectAttempt is called before each reconnection attempt.
	ReconnectAttempt(attempt int, delay time.Duration)

	// EnvelopeReceived is called when an envelope is received from Slack.
	EnvelopeReceived(envType string)

	// EnvelopeAcked is called when an envelope acknowledgment is sent.
	EnvelopeAcked(envType string, latency time.Duration)

	// HandlerStarted is called when a handler begins processing.
	HandlerStarted(envType string)

	// HandlerCompleted is called when a handler finishes processing.
	HandlerCompleted(envType string, duration time.Duration, err error)

	// HandlerPanic is called when a handler panics.
	HandlerPanic(envType string, recovered any)

	// WriteQueueDepth is called periodically with the current write queue depth.
	WriteQueueDepth(depth int)

	// PingSent is called when a ping is sent to the server.
	PingSent()

	// PongReceived is called when a pong is received from the server.
	PongReceived(latency time.Duration)

	// PongTimeout is called when a pong is not received in time.
	PongTimeout()
}

// NoopMetrics is a no-op implementation of MetricsHook.
// Use this as a default when no metrics collection is needed.
type NoopMetrics struct{}

var _ MetricsHook = (*NoopMetrics)(nil)

func (n *NoopMetrics) ConnectionOpened(connID string)                       {}
func (n *NoopMetrics) ConnectionClosed(connID string, duration time.Duration) {}
func (n *NoopMetrics) ReconnectAttempt(attempt int, delay time.Duration)    {}
func (n *NoopMetrics) EnvelopeReceived(envType string)                      {}
func (n *NoopMetrics) EnvelopeAcked(envType string, latency time.Duration)  {}
func (n *NoopMetrics) HandlerStarted(envType string)                        {}
func (n *NoopMetrics) HandlerCompleted(envType string, duration time.Duration, err error) {}
func (n *NoopMetrics) HandlerPanic(envType string, recovered any)           {}
func (n *NoopMetrics) WriteQueueDepth(depth int)                            {}
func (n *NoopMetrics) PingSent()                                            {}
func (n *NoopMetrics) PongReceived(latency time.Duration)                   {}
func (n *NoopMetrics) PongTimeout()                                         {}
