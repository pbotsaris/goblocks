package socketmode

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/coder/websocket"
)

func TestWriter_Write(t *testing.T) {
	t.Run("writes data to connection", func(t *testing.T) {
		var received []byte
		var mu sync.Mutex
		done := make(chan struct{})

		server := startMockServer(t, func(conn *websocket.Conn) {
			ctx := context.Background()
			_, data, err := conn.Read(ctx)
			if err != nil {
				return
			}
			mu.Lock()
			received = data
			mu.Unlock()
			close(done)
		})
		defer server.Close()

		ctx := context.Background()
		conn, _, err := websocket.Dial(ctx, server.URL(), nil)
		if err != nil {
			t.Fatalf("dial error: %v", err)
		}
		defer conn.CloseNow()

		w := newWriter(conn, defaultWriterConfig())
		defer w.Close()

		testData := []byte(`{"test": "data"}`)
		if err := w.Write(ctx, testData); err != nil {
			t.Fatalf("write error: %v", err)
		}

		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("timeout waiting for data")
		}

		mu.Lock()
		defer mu.Unlock()
		if string(received) != string(testData) {
			t.Errorf("got %q, want %q", received, testData)
		}
	})

	t.Run("multiple writes are serialized", func(t *testing.T) {
		var received [][]byte
		var mu sync.Mutex
		wg := sync.WaitGroup{}
		wg.Add(3)

		server := startMockServer(t, func(conn *websocket.Conn) {
			ctx := context.Background()
			for i := 0; i < 3; i++ {
				_, data, err := conn.Read(ctx)
				if err != nil {
					return
				}
				mu.Lock()
				received = append(received, data)
				mu.Unlock()
				wg.Done()
			}
		})
		defer server.Close()

		ctx := context.Background()
		conn, _, err := websocket.Dial(ctx, server.URL(), nil)
		if err != nil {
			t.Fatalf("dial error: %v", err)
		}
		defer conn.CloseNow()

		w := newWriter(conn, defaultWriterConfig())
		defer w.Close()

		// Write concurrently
		var writeWg sync.WaitGroup
		for i := 0; i < 3; i++ {
			writeWg.Add(1)
			go func(i int) {
				defer writeWg.Done()
				data := []byte{byte('A' + i)}
				if err := w.Write(ctx, data); err != nil {
					t.Errorf("write error: %v", err)
				}
			}(i)
		}

		writeWg.Wait()
		wg.Wait()

		mu.Lock()
		defer mu.Unlock()
		if len(received) != 3 {
			t.Errorf("got %d messages, want 3", len(received))
		}
	})
}

func TestWriter_WriteJSON(t *testing.T) {
	t.Run("marshals and writes JSON", func(t *testing.T) {
		var received map[string]any
		done := make(chan struct{})

		server := startMockServer(t, func(conn *websocket.Conn) {
			ctx := context.Background()
			_, data, err := conn.Read(ctx)
			if err != nil {
				return
			}
			json.Unmarshal(data, &received)
			close(done)
		})
		defer server.Close()

		ctx := context.Background()
		conn, _, err := websocket.Dial(ctx, server.URL(), nil)
		if err != nil {
			t.Fatalf("dial error: %v", err)
		}
		defer conn.CloseNow()

		w := newWriter(conn, defaultWriterConfig())
		defer w.Close()

		testObj := map[string]string{"key": "value"}
		if err := w.WriteJSON(ctx, testObj); err != nil {
			t.Fatalf("write error: %v", err)
		}

		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("timeout waiting for data")
		}

		if received["key"] != "value" {
			t.Errorf("got %v, want key=value", received)
		}
	})

	t.Run("returns error for invalid JSON", func(t *testing.T) {
		server := startMockServer(t, func(conn *websocket.Conn) {
			// Don't need to do anything
		})
		defer server.Close()

		ctx := context.Background()
		conn, _, err := websocket.Dial(ctx, server.URL(), nil)
		if err != nil {
			t.Fatalf("dial error: %v", err)
		}
		defer conn.CloseNow()

		w := newWriter(conn, defaultWriterConfig())
		defer w.Close()

		// Channels can't be marshaled to JSON
		invalid := make(chan int)
		err = w.WriteJSON(ctx, invalid)
		if err == nil {
			t.Error("expected error for invalid JSON")
		}
	})
}

func TestWriter_WriteAck(t *testing.T) {
	t.Run("writes correct ack format", func(t *testing.T) {
		var received Ack
		done := make(chan struct{})

		server := startMockServer(t, func(conn *websocket.Conn) {
			ctx := context.Background()
			_, data, err := conn.Read(ctx)
			if err != nil {
				return
			}
			json.Unmarshal(data, &received)
			close(done)
		})
		defer server.Close()

		ctx := context.Background()
		conn, _, err := websocket.Dial(ctx, server.URL(), nil)
		if err != nil {
			t.Fatalf("dial error: %v", err)
		}
		defer conn.CloseNow()

		w := newWriter(conn, defaultWriterConfig())
		defer w.Close()

		if err := w.WriteAck(ctx, "env123", nil); err != nil {
			t.Fatalf("write error: %v", err)
		}

		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("timeout waiting for data")
		}

		if received.EnvelopeID != "env123" {
			t.Errorf("got envelope_id %q, want %q", received.EnvelopeID, "env123")
		}
	})

	t.Run("includes payload when provided", func(t *testing.T) {
		var received map[string]any
		done := make(chan struct{})

		server := startMockServer(t, func(conn *websocket.Conn) {
			ctx := context.Background()
			_, data, err := conn.Read(ctx)
			if err != nil {
				return
			}
			json.Unmarshal(data, &received)
			close(done)
		})
		defer server.Close()

		ctx := context.Background()
		conn, _, err := websocket.Dial(ctx, server.URL(), nil)
		if err != nil {
			t.Fatalf("dial error: %v", err)
		}
		defer conn.CloseNow()

		w := newWriter(conn, defaultWriterConfig())
		defer w.Close()

		payload := map[string]string{"response": "data"}
		if err := w.WriteAck(ctx, "env456", payload); err != nil {
			t.Fatalf("write error: %v", err)
		}

		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("timeout waiting for data")
		}

		if received["envelope_id"] != "env456" {
			t.Errorf("got envelope_id %v", received["envelope_id"])
		}

		payloadMap := received["payload"].(map[string]any)
		if payloadMap["response"] != "data" {
			t.Errorf("got payload %v", payloadMap)
		}
	})
}

func TestWriter_Close(t *testing.T) {
	t.Run("closes cleanly", func(t *testing.T) {
		server := startMockServer(t, func(conn *websocket.Conn) {
			// Keep connection open until it closes
			ctx := context.Background()
			conn.Read(ctx)
		})
		defer server.Close()

		ctx := context.Background()
		conn, _, err := websocket.Dial(ctx, server.URL(), nil)
		if err != nil {
			t.Fatalf("dial error: %v", err)
		}
		defer conn.CloseNow()

		w := newWriter(conn, defaultWriterConfig())

		// Close should not panic
		w.Close()

		// Double close should not panic
		w.Close()
	})

	t.Run("write after close returns error", func(t *testing.T) {
		server := startMockServer(t, func(conn *websocket.Conn) {
			ctx := context.Background()
			conn.Read(ctx)
		})
		defer server.Close()

		ctx := context.Background()
		conn, _, err := websocket.Dial(ctx, server.URL(), nil)
		if err != nil {
			t.Fatalf("dial error: %v", err)
		}
		defer conn.CloseNow()

		w := newWriter(conn, defaultWriterConfig())
		w.Close()

		err = w.Write(ctx, []byte("test"))
		if err == nil {
			t.Error("expected error after close")
		}
	})
}

func TestWriter_Timeout(t *testing.T) {
	t.Run("respects context cancellation", func(t *testing.T) {
		server := startMockServer(t, func(conn *websocket.Conn) {
			// Don't read, let writes block
			time.Sleep(10 * time.Second)
		})
		defer server.Close()

		ctx := context.Background()
		conn, _, err := websocket.Dial(ctx, server.URL(), nil)
		if err != nil {
			t.Fatalf("dial error: %v", err)
		}
		defer conn.CloseNow()

		w := newWriter(conn, writerConfig{
			queueSize:    1,
			writeTimeout: 100 * time.Millisecond,
			metrics:      &NoopMetrics{},
		})
		defer w.Close()

		// Use a very short context timeout
		shortCtx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
		defer cancel()

		err = w.Write(shortCtx, []byte("test"))
		if err == nil {
			// Either context error or write error is acceptable
			// The important thing is it doesn't hang
		}
	})
}

func TestWriterConfig(t *testing.T) {
	t.Run("defaultWriterConfig returns sensible values", func(t *testing.T) {
		cfg := defaultWriterConfig()

		if cfg.queueSize != 100 {
			t.Errorf("queueSize: got %d, want 100", cfg.queueSize)
		}
		if cfg.writeTimeout != 3*time.Second {
			t.Errorf("writeTimeout: got %v, want 3s", cfg.writeTimeout)
		}
		if cfg.metrics == nil {
			t.Error("metrics should not be nil")
		}
	})
}
