# socketmode

Production-ready Socket Mode client for Slack, with type-safe responses.

[Back to main documentation](../README.md)

## Quick Start

```go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/pbotsaris/goblocks/blocks"
    "github.com/pbotsaris/goblocks/socketmode"
)

func main() {
    client := socketmode.New(os.Getenv("SLACK_APP_TOKEN"))

    client.OnSlashCommand(func(ctx context.Context, env *socketmode.Envelope) socketmode.Response {
        msg := blocks.NewBuilder().
            AddSection(blocks.MustMarkdown("Hello from */mycommand*!")).
            MustToMessage("Hello!")
        return socketmode.RespondWithMessage(msg)
    })

    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop()

    if err := client.Run(ctx); err != nil {
        log.Fatal(err)
    }
}
```

## Handler Registration

```go
// Convenience methods for common event types
client.OnSlashCommand(handler)  // slash_commands
client.OnInteractive(handler)   // interactive (buttons, modals, etc.)
client.OnEventsAPI(handler)     // events_api

// Generic handler for any event type
client.On("custom_event", handler)
```

The handler signature:

```go
type EventHandler func(ctx context.Context, envelope *Envelope) Response
```

The `Envelope` contains:

```go
type Envelope struct {
    EnvelopeID             string          // Unique ID for acknowledgment
    Type                   string          // Event type (events_api, interactive, etc.)
    Payload                json.RawMessage // Raw event payload
    AcceptsResponsePayload bool            // Whether response can include data
    RetryAttempt           int             // Retry attempt number (0 = first try)
    RetryReason            string          // Why this is a retry
}
```

## Type-Safe Responses

Response builders integrate with the [blocks](../blocks/README.md) package:

```go
// Empty response (ack only)
socketmode.NoResponse()

// Message response (slash commands)
msg := blocks.NewBuilder().
    AddSection(blocks.MustMarkdown("*Result:* Success")).
    MustToMessage("Result")
socketmode.RespondWithMessage(msg)

// Quick message from blocks
socketmode.RespondWithBlocks([]blocks.Block{section, divider})

// Modal responses (view submissions)
socketmode.RespondWithModalUpdate(modal)  // Replace current modal
socketmode.RespondWithModalPush(modal)    // Push new modal onto stack
socketmode.RespondWithModalClear()        // Close all modals
socketmode.RespondWithErrors(map[string]string{
    "email_block": "Invalid email address",
})

// Dynamic options (external selects)
socketmode.RespondWithOptions([]blocks.Option{opt1, opt2})
socketmode.RespondWithOptionGroups([]blocks.OptionGroup{group1})
```

## Client Options

```go
client := socketmode.New(appToken,
    socketmode.WithLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil))),
    socketmode.WithMetrics(&myMetrics{}),
    socketmode.WithMaxConcurrency(20),
    socketmode.WithHandlerTimeout(60 * time.Second),
    socketmode.WithHelloTimeout(30 * time.Second),
    socketmode.WithPingInterval(5 * time.Second),
    socketmode.WithPongTimeout(10 * time.Second),
    socketmode.WithHTTPClient(&http.Client{Timeout: 30 * time.Second}),
)
```

| Option | Default | Description |
|--------|---------|-------------|
| `WithLogger` | `slog.Default()` | Structured logger for debug/error output |
| `WithMetrics` | `NoopMetrics{}` | Metrics hook for observability |
| `WithMaxConcurrency` | `10` | Max concurrent handler goroutines |
| `WithHandlerTimeout` | `30s` | Timeout for handler execution |
| `WithHelloTimeout` | `30s` | Timeout waiting for hello message |
| `WithPingInterval` | `5s` | Interval between ping messages |
| `WithPongTimeout` | `10s` | Timeout waiting for pong response |
| `WithHTTPClient` | Default client | HTTP client for API calls |

## Metrics & Observability

Implement `MetricsHook` to collect metrics:

```go
type MetricsHook interface {
    ConnectionOpened(connID string)
    ConnectionClosed(connID string, duration time.Duration)
    ReconnectAttempt(attempt int, delay time.Duration)
    EnvelopeReceived(envType string)
    EnvelopeAcked(envType string, latency time.Duration)
    HandlerStarted(envType string)
    HandlerCompleted(envType string, duration time.Duration, err error)
    HandlerPanic(envType string, recovered any)
    WriteQueueDepth(depth int)
    PingSent()
    PongReceived(latency time.Duration)
    PongTimeout()
}
```

Example with Prometheus:

```go
type PrometheusMetrics struct {
    envelopesReceived *prometheus.CounterVec
    ackLatency        *prometheus.HistogramVec
    handlerDuration   *prometheus.HistogramVec
}

func (m *PrometheusMetrics) EnvelopeReceived(envType string) {
    m.envelopesReceived.WithLabelValues(envType).Inc()
}

func (m *PrometheusMetrics) EnvelopeAcked(envType string, latency time.Duration) {
    m.ackLatency.WithLabelValues(envType).Observe(latency.Seconds())
}

func (m *PrometheusMetrics) HandlerCompleted(envType string, duration time.Duration, err error) {
    m.handlerDuration.WithLabelValues(envType).Observe(duration.Seconds())
}
```

## Running with an HTTP Server

Run Socket Mode alongside an HTTP server using `errgroup`:

```go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"

    "golang.org/x/sync/errgroup"
    "github.com/pbotsaris/goblocks/socketmode"
)

func main() {
    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop()

    client := socketmode.New(os.Getenv("SLACK_APP_TOKEN"))
    client.OnSlashCommand(handleSlashCommand)

    mux := http.NewServeMux()
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    })

    server := &http.Server{Addr: ":8080", Handler: mux}

    g, ctx := errgroup.WithContext(ctx)

    g.Go(func() error {
        return client.Run(ctx)
    })

    g.Go(func() error {
        return server.ListenAndServe()
    })

    g.Go(func() error {
        <-ctx.Done()
        return server.Shutdown(context.Background())
    })

    if err := g.Wait(); err != nil && err != http.ErrServerClosed {
        log.Fatal(err)
    }
}
```

## Error Handling

The client classifies errors automatically:

**Permanent errors** (stops reconnection):
- Invalid authentication (`invalid_auth`)
- Token revoked (`token_revoked`)
- App uninstalled (`app_uninstalled`)
- Socket Mode disabled (`link_disabled`)
- HTTP 401, 403

**Retryable errors** (triggers reconnection with backoff):
- Network timeouts
- Connection refused
- HTTP 429, 5xx
- Rate limiting

Exponential backoff:
- Base delay: 1 second
- Max delay: 30 seconds
- Jitter: 0-1 second
- Resets after 60 seconds of stable connection

## Panic Recovery

Handler panics are recovered automatically:

```go
client.OnSlashCommand(func(ctx context.Context, env *socketmode.Envelope) socketmode.Response {
    panic("oops") // Recovered! Connection continues.
})
```

## Complete Example: Modal Workflow

```go
package main

import (
    "context"
    "encoding/json"
    "log/slog"
    "os"
    "os/signal"
    "syscall"

    "github.com/pbotsaris/goblocks/socketmode"
)

func main() {
    logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

    client := socketmode.New(
        os.Getenv("SLACK_APP_TOKEN"),
        socketmode.WithLogger(logger),
    )

    client.OnInteractive(func(ctx context.Context, env *socketmode.Envelope) socketmode.Response {
        var payload struct {
            Type string `json:"type"`
        }
        if err := json.Unmarshal(env.Payload, &payload); err != nil {
            return socketmode.NoResponse()
        }

        switch payload.Type {
        case "view_submission":
            return handleViewSubmission(env)
        default:
            return socketmode.NoResponse()
        }
    })

    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop()

    if err := client.Run(ctx); err != nil {
        logger.Error("client stopped", "error", err)
    }
}

func handleViewSubmission(env *socketmode.Envelope) socketmode.Response {
    var submission struct {
        View struct {
            State struct {
                Values map[string]map[string]struct {
                    Value string `json:"value"`
                } `json:"values"`
            } `json:"state"`
        } `json:"view"`
    }

    if err := json.Unmarshal(env.Payload, &submission); err != nil {
        return socketmode.NoResponse()
    }

    email := submission.View.State.Values["email_block"]["email_input"].Value
    if email == "" {
        return socketmode.RespondWithErrors(map[string]string{
            "email_block": "Email is required",
        })
    }

    return socketmode.RespondWithModalClear()
}
```
