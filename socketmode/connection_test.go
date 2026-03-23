package socketmode

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/coder/websocket"
)

func TestConnection_HelloHandshake(t *testing.T) {
	t.Run("waits for and parses hello message", func(t *testing.T) {
		server := startMockServer(t, func(conn *websocket.Conn) {
			sendHello(t, conn, "A123", 1)
			// Keep connection open
			time.Sleep(100 * time.Millisecond)
		})
		defer server.Close()

		ctx := context.Background()
		metrics := &mockMetrics{}

		cfg := connectionConfig{
			id:             "test-conn",
			url:            server.URL(),
			helloTimeout:   5 * time.Second,
			writeQueueSize: 10,
			writeTimeout:   3 * time.Second,
			logger:         slog.Default(),
			metrics:        metrics,
		}

		conn, err := dial(ctx, cfg)
		if err != nil {
			t.Fatalf("dial error: %v", err)
		}
		defer conn.Close()

		// Check hello info
		hello := conn.HelloInfo()
		if hello == nil {
			t.Fatal("expected hello info")
		}
		if hello.ConnectionInfo.AppID != "A123" {
			t.Errorf("got app_id %q, want %q", hello.ConnectionInfo.AppID, "A123")
		}
		if hello.NumConnections != 1 {
			t.Errorf("got num_connections %d, want 1", hello.NumConnections)
		}

		// Check metrics
		if metrics.connectionsOpened != 1 {
			t.Errorf("got %d connections opened, want 1", metrics.connectionsOpened)
		}
	})

	t.Run("times out if no hello received", func(t *testing.T) {
		server := startMockServer(t, func(conn *websocket.Conn) {
			// Don't send hello, just wait
			time.Sleep(5 * time.Second)
		})
		defer server.Close()

		ctx := context.Background()

		cfg := connectionConfig{
			id:             "test-conn",
			url:            server.URL(),
			helloTimeout:   100 * time.Millisecond, // Very short
			writeQueueSize: 10,
			writeTimeout:   3 * time.Second,
			logger:         slog.Default(),
			metrics:        &NoopMetrics{},
		}

		_, err := dial(ctx, cfg)
		if err == nil {
			t.Fatal("expected timeout error")
		}
	})

	t.Run("connection ID is set", func(t *testing.T) {
		server := startMockServer(t, func(conn *websocket.Conn) {
			sendHello(t, conn, "A123", 1)
			time.Sleep(100 * time.Millisecond)
		})
		defer server.Close()

		ctx := context.Background()

		cfg := connectionConfig{
			id:             "my-conn-id",
			url:            server.URL(),
			helloTimeout:   5 * time.Second,
			writeQueueSize: 10,
			writeTimeout:   3 * time.Second,
			logger:         slog.Default(),
			metrics:        &NoopMetrics{},
		}

		conn, err := dial(ctx, cfg)
		if err != nil {
			t.Fatalf("dial error: %v", err)
		}
		defer conn.Close()

		if conn.ID() != "my-conn-id" {
			t.Errorf("got id %q, want %q", conn.ID(), "my-conn-id")
		}
	})
}

func TestConnection_Read(t *testing.T) {
	t.Run("returns message data", func(t *testing.T) {
		server := startMockServer(t, func(conn *websocket.Conn) {
			sendHello(t, conn, "A123", 1)

			// Send an envelope
			env := Envelope{
				EnvelopeID: "env123",
				Type:       "events_api",
			}
			sendEnvelope(t, conn, env)
		})
		defer server.Close()

		ctx := context.Background()

		cfg := connectionConfig{
			id:             "test-conn",
			url:            server.URL(),
			helloTimeout:   5 * time.Second,
			writeQueueSize: 10,
			writeTimeout:   3 * time.Second,
			logger:         slog.Default(),
			metrics:        &NoopMetrics{},
		}

		conn, err := dial(ctx, cfg)
		if err != nil {
			t.Fatalf("dial error: %v", err)
		}
		defer conn.Close()

		// Read the envelope
		data, err := conn.Read(ctx)
		if err != nil {
			t.Fatalf("read error: %v", err)
		}

		var env Envelope
		mustUnmarshal(t, data, &env)

		if env.EnvelopeID != "env123" {
			t.Errorf("got envelope_id %q, want %q", env.EnvelopeID, "env123")
		}
	})
}

func TestConnection_Close(t *testing.T) {
	t.Run("closes cleanly", func(t *testing.T) {
		server := startMockServer(t, func(conn *websocket.Conn) {
			sendHello(t, conn, "A123", 1)
			// Wait for close
			ctx := context.Background()
			conn.Read(ctx) // Will return error when closed
		})
		defer server.Close()

		ctx := context.Background()
		metrics := &mockMetrics{}

		cfg := connectionConfig{
			id:             "test-conn",
			url:            server.URL(),
			helloTimeout:   5 * time.Second,
			writeQueueSize: 10,
			writeTimeout:   3 * time.Second,
			logger:         slog.Default(),
			metrics:        metrics,
		}

		conn, err := dial(ctx, cfg)
		if err != nil {
			t.Fatalf("dial error: %v", err)
		}

		conn.Close()

		// Check Done() is closed
		select {
		case <-conn.Done():
			// Good
		case <-time.After(1 * time.Second):
			t.Error("Done() not closed after Close()")
		}

		// Check metrics
		if metrics.connectionsClosed != 1 {
			t.Errorf("got %d connections closed, want 1", metrics.connectionsClosed)
		}
	})

	t.Run("double close is safe", func(t *testing.T) {
		server := startMockServer(t, func(conn *websocket.Conn) {
			sendHello(t, conn, "A123", 1)
			time.Sleep(1 * time.Second)
		})
		defer server.Close()

		ctx := context.Background()

		cfg := connectionConfig{
			id:             "test-conn",
			url:            server.URL(),
			helloTimeout:   5 * time.Second,
			writeQueueSize: 10,
			writeTimeout:   3 * time.Second,
			logger:         slog.Default(),
			metrics:        &NoopMetrics{},
		}

		conn, err := dial(ctx, cfg)
		if err != nil {
			t.Fatalf("dial error: %v", err)
		}

		// Should not panic
		conn.Close()
		conn.Close()
	})
}

func TestConnection_Writer(t *testing.T) {
	t.Run("returns working writer", func(t *testing.T) {
		ackReceived := make(chan Ack, 1)

		server := startMockServer(t, func(conn *websocket.Conn) {
			sendHello(t, conn, "A123", 1)
			ack := readAck(t, conn)
			ackReceived <- ack
		})
		defer server.Close()

		ctx := context.Background()

		cfg := connectionConfig{
			id:             "test-conn",
			url:            server.URL(),
			helloTimeout:   5 * time.Second,
			writeQueueSize: 10,
			writeTimeout:   3 * time.Second,
			logger:         slog.Default(),
			metrics:        &NoopMetrics{},
		}

		conn, err := dial(ctx, cfg)
		if err != nil {
			t.Fatalf("dial error: %v", err)
		}
		defer conn.Close()

		writer := conn.Writer()
		if writer == nil {
			t.Fatal("expected writer")
		}

		// Write an ack
		err = writer.WriteAck(ctx, "test-env", nil)
		if err != nil {
			t.Fatalf("write error: %v", err)
		}

		select {
		case ack := <-ackReceived:
			if ack.EnvelopeID != "test-env" {
				t.Errorf("got envelope_id %q, want %q", ack.EnvelopeID, "test-env")
			}
		case <-time.After(5 * time.Second):
			t.Fatal("timeout waiting for ack")
		}
	})
}

func TestConnection_ContextCancellation(t *testing.T) {
	t.Run("cancels dial on context cancellation", func(t *testing.T) {
		server := startMockServer(t, func(conn *websocket.Conn) {
			// Never send hello
			time.Sleep(10 * time.Second)
		})
		defer server.Close()

		ctx, cancel := context.WithCancel(context.Background())

		cfg := connectionConfig{
			id:             "test-conn",
			url:            server.URL(),
			helloTimeout:   10 * time.Second,
			writeQueueSize: 10,
			writeTimeout:   3 * time.Second,
			logger:         slog.Default(),
			metrics:        &NoopMetrics{},
		}

		// Cancel after a short delay
		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()

		_, err := dial(ctx, cfg)
		if err == nil {
			t.Fatal("expected error on context cancellation")
		}
	})
}
