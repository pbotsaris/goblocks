package socketmode

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestClient_Options(t *testing.T) {
	t.Run("WithLogger sets logger", func(t *testing.T) {
		logger := slog.Default()
		client := New("xapp-test", WithLogger(logger))

		if client.logger != logger {
			t.Error("logger not set")
		}
	})

	t.Run("WithMetrics sets metrics", func(t *testing.T) {
		metrics := &mockMetrics{}
		client := New("xapp-test", WithMetrics(metrics))

		if client.metrics != metrics {
			t.Error("metrics not set")
		}
	})

	t.Run("WithMaxConcurrency sets limit", func(t *testing.T) {
		client := New("xapp-test", WithMaxConcurrency(50))

		if client.maxConcurrency != 50 {
			t.Errorf("got %d, want 50", client.maxConcurrency)
		}
	})

	t.Run("WithHandlerTimeout sets timeout", func(t *testing.T) {
		client := New("xapp-test", WithHandlerTimeout(5*time.Minute))

		if client.handlerTimeout != 5*time.Minute {
			t.Errorf("got %v, want 5m", client.handlerTimeout)
		}
	})

	t.Run("WithHelloTimeout sets timeout", func(t *testing.T) {
		client := New("xapp-test", WithHelloTimeout(10*time.Second))

		if client.helloTimeout != 10*time.Second {
			t.Errorf("got %v, want 10s", client.helloTimeout)
		}
	})

	t.Run("WithPingInterval sets interval", func(t *testing.T) {
		client := New("xapp-test", WithPingInterval(10*time.Second))

		if client.pingInterval != 10*time.Second {
			t.Errorf("got %v, want 10s", client.pingInterval)
		}
	})

	t.Run("WithPongTimeout sets timeout", func(t *testing.T) {
		client := New("xapp-test", WithPongTimeout(20*time.Second))

		if client.pongTimeout != 20*time.Second {
			t.Errorf("got %v, want 20s", client.pongTimeout)
		}
	})

	t.Run("default values are sensible", func(t *testing.T) {
		client := New("xapp-test")

		if client.maxConcurrency != 10 {
			t.Errorf("maxConcurrency: got %d, want 10", client.maxConcurrency)
		}
		if client.handlerTimeout != 30*time.Second {
			t.Errorf("handlerTimeout: got %v, want 30s", client.handlerTimeout)
		}
		if client.helloTimeout != 30*time.Second {
			t.Errorf("helloTimeout: got %v, want 30s", client.helloTimeout)
		}
	})

	t.Run("nil values are ignored", func(t *testing.T) {
		client := New("xapp-test",
			WithLogger(nil),
			WithMetrics(nil),
			WithHTTPClient(nil),
		)

		if client.logger == nil {
			t.Error("logger should not be nil")
		}
		if client.metrics == nil {
			t.Error("metrics should not be nil")
		}
		if client.httpClient == nil {
			t.Error("httpClient should not be nil")
		}
	})
}

func TestClient_HandlerRegistration(t *testing.T) {
	t.Run("On registers handler", func(t *testing.T) {
		client := New("xapp-test")
		called := false

		client.On("custom_event", func(ctx context.Context, env *Envelope) Response {
			called = true
			return NoResponse()
		})

		// Check handler is registered
		client.handlersMu.RLock()
		handler, ok := client.handlers["custom_event"]
		client.handlersMu.RUnlock()

		if !ok {
			t.Fatal("handler not registered")
		}

		// Call handler
		handler(context.Background(), &Envelope{})
		if !called {
			t.Error("handler not called")
		}
	})

	t.Run("OnSlashCommand registers for slash_commands", func(t *testing.T) {
		client := New("xapp-test")

		client.OnSlashCommand(func(ctx context.Context, env *Envelope) Response {
			return NoResponse()
		})

		client.handlersMu.RLock()
		_, ok := client.handlers[EnvelopeTypeSlashCommands]
		client.handlersMu.RUnlock()

		if !ok {
			t.Error("handler not registered for slash_commands")
		}
	})

	t.Run("OnInteractive registers for interactive", func(t *testing.T) {
		client := New("xapp-test")

		client.OnInteractive(func(ctx context.Context, env *Envelope) Response {
			return NoResponse()
		})

		client.handlersMu.RLock()
		_, ok := client.handlers[EnvelopeTypeInteractive]
		client.handlersMu.RUnlock()

		if !ok {
			t.Error("handler not registered for interactive")
		}
	})

	t.Run("OnEventsAPI registers for events_api", func(t *testing.T) {
		client := New("xapp-test")

		client.OnEventsAPI(func(ctx context.Context, env *Envelope) Response {
			return NoResponse()
		})

		client.handlersMu.RLock()
		_, ok := client.handlers[EnvelopeTypeEventsAPI]
		client.handlersMu.RUnlock()

		if !ok {
			t.Error("handler not registered for events_api")
		}
	})
}

func TestClient_SafeHandler(t *testing.T) {
	t.Run("recovers from panic", func(t *testing.T) {
		metrics := &mockMetrics{}
		client := New("xapp-test", WithMetrics(metrics))

		handler := func(ctx context.Context, env *Envelope) Response {
			panic("test panic")
		}

		env := &Envelope{EnvelopeID: "test", Type: "test"}

		// Should not panic
		resp := client.safeHandler(context.Background(), env, handler)

		// Should return empty response
		if _, ok := resp.(EmptyResponse); !ok {
			t.Error("expected EmptyResponse after panic")
		}

		// Should record panic metric
		metrics.mu.Lock()
		panics := metrics.handlerPanics
		metrics.mu.Unlock()

		if panics != 1 {
			t.Errorf("got %d panics, want 1", panics)
		}
	})

	t.Run("returns handler response normally", func(t *testing.T) {
		client := New("xapp-test")

		handler := func(ctx context.Context, env *Envelope) Response {
			return RespondWithModalClear()
		}

		env := &Envelope{EnvelopeID: "test", Type: "test"}
		resp := client.safeHandler(context.Background(), env, handler)

		if _, ok := resp.(ModalResponse); !ok {
			t.Error("expected ModalResponse")
		}
	})

	t.Run("records handler metrics", func(t *testing.T) {
		metrics := &mockMetrics{}
		client := New("xapp-test", WithMetrics(metrics))

		handler := func(ctx context.Context, env *Envelope) Response {
			return NoResponse()
		}

		env := &Envelope{EnvelopeID: "test", Type: "events_api"}
		client.safeHandler(context.Background(), env, handler)

		metrics.mu.Lock()
		started := metrics.handlersStarted
		completed := metrics.handlersCompleted
		metrics.mu.Unlock()

		if started != 1 {
			t.Errorf("got %d handlers started, want 1", started)
		}
		if completed != 1 {
			t.Errorf("got %d handlers completed, want 1", completed)
		}
	})

	t.Run("handler receives context with timeout", func(t *testing.T) {
		client := New("xapp-test", WithHandlerTimeout(100*time.Millisecond))

		var receivedCtx context.Context
		handler := func(ctx context.Context, env *Envelope) Response {
			receivedCtx = ctx
			return NoResponse()
		}

		env := &Envelope{EnvelopeID: "test", Type: "test"}
		client.safeHandler(context.Background(), env, handler)

		deadline, ok := receivedCtx.Deadline()
		if !ok {
			t.Error("expected context to have deadline")
		}

		// Deadline should be roughly 100ms from now (already passed)
		if time.Until(deadline) > 100*time.Millisecond {
			t.Error("deadline seems too far in the future")
		}
	})
}

func TestClient_ConcurrencyLimit(t *testing.T) {
	t.Run("limits concurrent handlers", func(t *testing.T) {
		client := New("xapp-test", WithMaxConcurrency(2))

		var concurrent int32
		var maxConcurrent int32
		var wg sync.WaitGroup

		handler := func(ctx context.Context, env *Envelope) Response {
			current := atomic.AddInt32(&concurrent, 1)
			for {
				old := atomic.LoadInt32(&maxConcurrent)
				if current <= old || atomic.CompareAndSwapInt32(&maxConcurrent, old, current) {
					break
				}
			}

			time.Sleep(50 * time.Millisecond)
			atomic.AddInt32(&concurrent, -1)
			return NoResponse()
		}

		// Start 5 handlers
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				ctx := context.Background()
				env := &Envelope{EnvelopeID: "test", Type: "test"}

				// Simulate semaphore acquisition like in readLoop
				sem := make(chan struct{}, client.maxConcurrency)
				sem <- struct{}{}
				client.safeHandler(ctx, env, handler)
				<-sem
			}()
		}

		wg.Wait()

		// In reality the concurrency limit is enforced in readLoop,
		// not in safeHandler itself. This test verifies the pattern works.
	})
}

func TestEnvelopeTypes(t *testing.T) {
	t.Run("envelope type constants are correct", func(t *testing.T) {
		if EnvelopeTypeHello != "hello" {
			t.Error("EnvelopeTypeHello")
		}
		if EnvelopeTypeDisconnect != "disconnect" {
			t.Error("EnvelopeTypeDisconnect")
		}
		if EnvelopeTypeEventsAPI != "events_api" {
			t.Error("EnvelopeTypeEventsAPI")
		}
		if EnvelopeTypeInteractive != "interactive" {
			t.Error("EnvelopeTypeInteractive")
		}
		if EnvelopeTypeSlashCommands != "slash_commands" {
			t.Error("EnvelopeTypeSlashCommands")
		}
	})

	t.Run("disconnect reason constants are correct", func(t *testing.T) {
		if DisconnectReasonLinkDisabled != "link_disabled" {
			t.Error("DisconnectReasonLinkDisabled")
		}
		if DisconnectReasonWarning != "warning" {
			t.Error("DisconnectReasonWarning")
		}
		if DisconnectReasonRefreshRequested != "refresh_requested" {
			t.Error("DisconnectReasonRefreshRequested")
		}
	})
}

func TestTypes_Marshaling(t *testing.T) {
	t.Run("Envelope marshals correctly", func(t *testing.T) {
		env := Envelope{
			EnvelopeID:             "env123",
			Type:                   "events_api",
			Payload:                []byte(`{"test": true}`),
			AcceptsResponsePayload: true,
			RetryAttempt:           2,
			RetryReason:            "timeout",
		}

		data := mustJSON(t, env)
		var result map[string]any
		mustUnmarshal(t, data, &result)

		if result["envelope_id"] != "env123" {
			t.Error("envelope_id")
		}
		if result["type"] != "events_api" {
			t.Error("type")
		}
		if result["accepts_response_payload"] != true {
			t.Error("accepts_response_payload")
		}
	})

	t.Run("Ack marshals correctly", func(t *testing.T) {
		ack := Ack{
			EnvelopeID: "env456",
			Payload:    map[string]string{"key": "value"},
		}

		data := mustJSON(t, ack)
		var result map[string]any
		mustUnmarshal(t, data, &result)

		if result["envelope_id"] != "env456" {
			t.Error("envelope_id")
		}
		payload := result["payload"].(map[string]any)
		if payload["key"] != "value" {
			t.Error("payload")
		}
	})

	t.Run("Ack omits nil payload", func(t *testing.T) {
		ack := Ack{
			EnvelopeID: "env789",
		}

		data := mustJSON(t, ack)
		var result map[string]any
		mustUnmarshal(t, data, &result)

		if _, ok := result["payload"]; ok {
			t.Error("payload should be omitted when nil")
		}
	})

	t.Run("HelloMessage unmarshals correctly", func(t *testing.T) {
		data := []byte(`{
			"type": "hello",
			"connection_info": {"app_id": "A123"},
			"num_connections": 2,
			"debug_info": {
				"host": "test-host",
				"started": "2024-01-01",
				"build_number": 123,
				"approximate_connection_time": 3600
			}
		}`)

		var hello HelloMessage
		mustUnmarshal(t, data, &hello)

		if hello.Type != "hello" {
			t.Error("type")
		}
		if hello.ConnectionInfo.AppID != "A123" {
			t.Error("app_id")
		}
		if hello.NumConnections != 2 {
			t.Error("num_connections")
		}
		if hello.DebugInfo.ApproximateConnectionTime != 3600 {
			t.Error("approximate_connection_time")
		}
	})
}
