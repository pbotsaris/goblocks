# goblocks

Go library for building Slack Block Kit UIs.

## Installation

```go
import "github.com/pbotsaris/goblocks/blocks"
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/pbotsaris/goblocks/blocks"
)

func main() {
    // Build a simple message
    message := blocks.NewBuilder().
        AddHeader("Welcome!").
        AddSection(blocks.MustMarkdown("Hello, *world*!")).
        AddDivider().
        AddActions([]blocks.ActionsElement{
            blocks.MustButton("Click me", blocks.WithActionID("btn_click")),
        }).
        MustToMessage("Welcome message")

    // Get JSON output
    data, _ := json.MarshalIndent(message, "", "  ")
    fmt.Println(string(data))
}
```

## Architecture

### Hierarchy

```
Surface (Message, Modal, HomeTab)
    |
    +-- Block (Section, Actions, Header, etc.)
            |
            +-- Element (Button, Select, DatePicker, etc.)
                    |
                    +-- Composition Object (PlainText, Markdown, Option, etc.)
```

### Type Safety

The library uses sealed marker interfaces to ensure elements are only used where valid:

```go
// These interfaces restrict where elements can be placed
type Block interface { ... }              // Top-level blocks
type SectionAccessory interface { ... }   // Valid in section blocks
type ActionsElement interface { ... }     // Valid in actions blocks
type ContextElement interface { ... }     // Valid in context blocks
type InputElement interface { ... }       // Valid in input blocks
```

For example, a `Button` implements both `SectionAccessory` and `ActionsElement`, so it can be used in both contexts. But an `EmailInput` only implements `InputElement`, so it can only be used in input blocks.

## Patterns

### Constructor Pattern

All types use a `New<Type>()` constructor that returns `(<Type>, error)`:

```go
button, err := blocks.NewButton("Click me")
if err != nil {
    // handle error
}
```

For convenience, `Must<Type>()` variants panic on error (useful in tests or static configurations):

```go
button := blocks.MustButton("Click me")
```

### Functional Options

All types support optional configuration via functional options:

```go
button, err := blocks.NewButton("Delete",
    blocks.WithActionID("delete_btn"),
    blocks.WithButtonStyle(blocks.ButtonStyleDanger),
    blocks.WithButtonConfirm(confirmDialog),
)
```

Options follow the naming convention `With<Type><Field>()`:

```go
// Button options
blocks.WithActionID("...")
blocks.WithValue("...")
blocks.WithButtonStyle(blocks.ButtonStylePrimary)
blocks.WithButtonConfirm(dialog)

// Section options
blocks.WithSectionBlockID("...")
blocks.WithSectionAccessory(button)
blocks.WithSectionFields(field1, field2)
```

### Builder Pattern

The `Builder` provides a fluent API for composing multiple blocks:

```go
builder := blocks.NewBuilder().
    AddHeader("Report").
    AddSection(blocks.MustMarkdown("*Summary*")).
    AddDivider().
    AddContext([]blocks.ContextElement{
        blocks.MustPlainText("Last updated: today"),
    })

// Convert to different surfaces
message, _ := builder.ToMessage("Report")
modal, _ := builder.ToModal("Report Modal")
homeTab, _ := builder.ToHomeTab()

// Or get raw blocks
blks, _ := builder.Build()

// Get JSON directly
jsonData, _ := builder.PrettyJSON()
```

## Surfaces

Surfaces are the top-level containers that hold blocks.

### Message

```go
// With builder
message := blocks.NewBuilder().
    AddHeader("Notification").
    AddSection(blocks.MustMarkdown("You have a new message")).
    MustToMessage("New notification")

// Direct construction
message, err := blocks.NewMessage("Fallback text", myBlocks,
    blocks.WithMessageThreadTS("1234567890.123456"),
    blocks.WithMessageMrkdwn(),
)
```

### Modal

```go
modal, err := blocks.NewModal("Settings", myBlocks,
    blocks.WithModalSubmit("Save"),
    blocks.WithModalClose("Cancel"),
    blocks.WithModalCallbackID("settings_modal"),
    blocks.WithModalPrivateMetadata(`{"user_id": "U123"}`),
    blocks.WithModalClearOnClose(),
    blocks.WithModalNotifyOnClose(),
)
```

### Home Tab

```go
homeTab, err := blocks.NewHomeTab(myBlocks,
    blocks.WithHomeTabCallbackID("home_view"),
    blocks.WithHomeTabPrivateMetadata("..."),
)
```

## Blocks

### Section

The most flexible block, supports text, fields, and accessories.

```go
// With text
section, _ := blocks.NewSection(blocks.MustMarkdown("*Bold* text"))

// With accessory
section, _ := blocks.NewSection(
    blocks.MustMarkdown("Choose an option:"),
    blocks.WithSectionAccessory(selectMenu),
)

// With fields (2-column layout)
section, _ := blocks.NewSectionWithFields([]blocks.TextObject{
    blocks.MustMarkdown("*Name:*\nJohn"),
    blocks.MustMarkdown("*Role:*\nAdmin"),
})
```

### Actions

Container for interactive elements (max 25 elements).

```go
actions, _ := blocks.NewActions([]blocks.ActionsElement{
    blocks.MustButton("Approve", blocks.WithButtonStyle(blocks.ButtonStylePrimary)),
    blocks.MustButton("Reject", blocks.WithButtonStyle(blocks.ButtonStyleDanger)),
    datePicker,
    selectMenu,
})
```

### Context

Displays secondary information (max 10 elements).

```go
context, _ := blocks.NewContext([]blocks.ContextElement{
    imageElement,
    blocks.MustMarkdown("Posted by <@U123>"),
    blocks.MustPlainText("2 hours ago"),
})
```

### Input

Collects user input in modals and messages.

```go
input, _ := blocks.NewInput("Email", blocks.NewEmailInput(),
    blocks.WithInputBlockID("email_input"),
    blocks.WithInputHint("Enter your work email"),
    blocks.WithInputOptional(),
)
```

### Header

Large text for section titles (max 150 characters).

```go
header, _ := blocks.NewHeader("Configuration")
```

### Divider

Visual separator between blocks.

```go
divider := blocks.NewDivider()
```

### Image

Displays an image.

```go
image, _ := blocks.NewImageBlock(
    "https://example.com/image.png",
    "Description of image",
    blocks.WithImageBlockTitle("My Image"),
)
```

### Video

Embedded video player.

```go
video, _ := blocks.NewVideo(
    "Video description",
    "My Video",
    "https://example.com/thumb.png",
    "https://example.com/video.mp4",
    blocks.WithVideoAuthorName("John Doe"),
    blocks.WithVideoProviderName("YouTube"),
)
```

### File

Displays remote file information (read-only, appears in retrieved messages).

```go
file, _ := blocks.NewFile("external_file_id")
```

### Rich Text

Structured formatted text with sections, lists, and quotes.

```go
section := blocks.NewRichTextSection([]blocks.RichTextSectionElement{
    blocks.NewRichTextText("Hello ", nil),
    blocks.NewRichTextText("world", blocks.NewRichTextStyle(true, false, false, false)), // bold
    blocks.NewRichTextEmoji("wave"),
})

list := blocks.NewRichTextList("bullet", []blocks.RichTextSection{
    blocks.NewRichTextSection([]blocks.RichTextSectionElement{
        blocks.NewRichTextText("First item", nil),
    }),
    blocks.NewRichTextSection([]blocks.RichTextSectionElement{
        blocks.NewRichTextText("Second item", nil),
    }),
})

richText, _ := blocks.NewRichText([]blocks.RichTextElement{section, list})
```

### Table

Structured tabular data (messages only).

```go
col1 := blocks.NewTableColumn("name", "Name")
col2 := blocks.NewTableColumn("value", "Value")

row := blocks.NewTableRow([]blocks.TableCell{
    blocks.NewTableCell("name", blocks.NewTableCellText("Item")),
    blocks.NewTableCell("value", blocks.NewTableCellText("$100")),
})

table, _ := blocks.NewTable([]blocks.TableColumn{col1, col2}, []blocks.TableRow{row})
```

### Plan

Collection of tasks (messages only).

```go
item1 := blocks.MustPlanItem("Setup environment", blocks.PlanItemStatusComplete)
item2 := blocks.MustPlanItem("Write tests", blocks.PlanItemStatusInProgress)
item3 := blocks.MustPlanItem("Deploy", blocks.PlanItemStatusPending)

section := blocks.MustPlanSection("Phase 1", []blocks.PlanItem{item1, item2, item3})
plan, _ := blocks.NewPlan("Project Plan", []blocks.PlanSection{section})
```

### Task Card

Single task display (messages only).

```go
source := blocks.MustURLSource("https://github.com/...", blocks.WithURLSourceTitle("GitHub"))

taskCard, _ := blocks.NewTaskCard("task_1", "Fix login bug", blocks.TaskCardStatusInProgress,
    blocks.WithTaskCardDescription("Users can't log in with SSO"),
    blocks.WithTaskCardSources([]blocks.URLSource{source}),
)
```

### Context Actions

AI/assistant feedback buttons (messages only).

```go
contextActions, _ := blocks.NewContextActions([]blocks.ContextActionsElement{
    blocks.NewFeedbackButtons(blocks.WithFeedbackButtonsActionID("feedback")),
    blocks.NewIconButton(blocks.NewIcon("copy"), blocks.WithIconButtonActionID("copy")),
})
```

### Markdown Block

Raw markdown content (messages only).

```go
md, _ := blocks.NewMarkdownBlock("# Heading\n\nParagraph with **bold** text.")
```

## Elements

### Button

```go
button, _ := blocks.NewButton("Click me",
    blocks.WithActionID("btn_click"),
    blocks.WithValue("clicked"),
    blocks.WithButtonStyle(blocks.ButtonStylePrimary), // or ButtonStyleDanger
    blocks.WithURL("https://example.com"),
    blocks.WithButtonConfirm(confirmDialog),
    blocks.WithAccessibilityLabel("Click this button"),
)
```

### Select Menus

```go
// Static select
option1 := blocks.MustOption("Option 1", "opt1")
option2 := blocks.MustOption("Option 2", "opt2")
staticSelect, _ := blocks.NewStaticSelect([]blocks.Option{option1, option2},
    blocks.WithStaticSelectActionID("select_action"),
    blocks.WithStaticSelectPlaceholder("Choose..."),
    blocks.WithStaticSelectInitial(option1),
)

// Users select
usersSelect := blocks.NewUsersSelect(
    blocks.WithUsersSelectActionID("user_select"),
    blocks.WithUsersSelectInitialUser("U123"),
)

// Conversations select
convoSelect := blocks.NewConversationsSelect(
    blocks.WithConversationsSelectActionID("convo_select"),
    blocks.WithConversationsSelectFilter(filter),
)

// Channels select
channelSelect := blocks.NewChannelsSelect(
    blocks.WithChannelsSelectActionID("channel_select"),
)

// External select (load options from your server)
externalSelect := blocks.NewExternalSelect(
    blocks.WithExternalSelectActionID("external_select"),
    blocks.WithExternalSelectMinQueryLength(3),
)
```

### Multi-Select Menus

All select types have multi-select variants:

```go
multiSelect, _ := blocks.NewMultiStaticSelect(options,
    blocks.WithMultiStaticSelectActionID("multi_select"),
    blocks.WithMultiStaticSelectMaxItems(5),
)
```

### Date & Time Pickers

```go
datePicker := blocks.NewDatePicker(
    blocks.WithDatePickerActionID("date_pick"),
    blocks.WithDatePickerInitialDate("2024-01-15"),
    blocks.WithDatePickerPlaceholder("Select date"),
)

timePicker := blocks.NewTimePicker(
    blocks.WithTimePickerActionID("time_pick"),
    blocks.WithTimePickerInitialTime("14:30"),
)

datetimePicker := blocks.NewDatetimePicker(
    blocks.WithDatetimePickerActionID("datetime_pick"),
    blocks.WithDatetimePickerInitialDateTime(1702656000), // Unix timestamp
)
```

### Checkboxes & Radio Buttons

```go
checkboxes, _ := blocks.NewCheckboxes(options,
    blocks.WithCheckboxesActionID("checkboxes"),
    blocks.WithCheckboxesInitialOptions(option1, option2),
)

radioButtons, _ := blocks.NewRadioButtons(options,
    blocks.WithRadioButtonsActionID("radio"),
    blocks.WithRadioButtonsInitialOption(option1),
)
```

### Input Elements

```go
// Plain text input
textInput := blocks.NewPlainTextInput(
    blocks.WithPlainTextInputActionID("text_input"),
    blocks.WithMultiline(),
    blocks.WithMinLength(10),
    blocks.WithMaxLength(500),
    blocks.WithPlainTextInputPlaceholder("Enter text..."),
)

// Email input
emailInput := blocks.NewEmailInput(
    blocks.WithEmailInputActionID("email"),
    blocks.WithEmailInputPlaceholder("you@example.com"),
)

// Number input
numberInput := blocks.NewNumberInput(true, // allow decimals
    blocks.WithNumberInputActionID("number"),
    blocks.WithNumberInputMinValue("0"),
    blocks.WithNumberInputMaxValue("100"),
)

// URL input
urlInput := blocks.NewURLInput(
    blocks.WithURLInputActionID("url"),
    blocks.WithURLInputPlaceholder("https://..."),
)

// File input
fileInput := blocks.NewFileInput(
    blocks.WithFileInputActionID("file"),
    blocks.WithFileInputFiletypes([]string{"pdf", "doc"}),
    blocks.WithFileInputMaxFiles(5),
)

// Rich text input
richTextInput := blocks.NewRichTextInput(
    blocks.WithRichTextInputActionID("rich_text"),
    blocks.WithRichTextInputPlaceholder("Write something..."),
)
```

### Overflow Menu

```go
overflow, _ := blocks.NewOverflow(options,
    blocks.WithOverflowActionID("overflow_menu"),
    blocks.WithOverflowConfirm(confirmDialog),
)
```

### Workflow Button

```go
trigger := blocks.MustTrigger("https://slack.com/shortcuts/...",
    blocks.WithInputParameters(
        blocks.MustInputParameter("user_id", "{{user.id}}"),
    ),
)

workflowBtn, _ := blocks.NewWorkflowButton("Start Workflow", blocks.NewWorkflow(trigger),
    blocks.WithWorkflowButtonStyle(blocks.ButtonStylePrimary),
)
```

### Image Element

For use in section accessories or context blocks:

```go
imageElem, _ := blocks.NewImageElement(
    "https://example.com/avatar.png",
    "User avatar",
)
```

## Composition Objects

### Text Objects

```go
// Plain text (emoji enabled by default)
plain := blocks.MustPlainText("Hello :wave:")

// Plain text without emoji
plain := blocks.MustPlainText("Hello :wave:", blocks.WithEmoji(false))

// Markdown
md := blocks.MustMarkdown("*Bold* and _italic_")

// Markdown without auto-linking
md := blocks.MustMarkdown("example.com", blocks.WithVerbatim(true))
```

### Option

```go
option := blocks.MustOption("Display Text", "value",
    blocks.WithDescription("Additional info"),
)
```

### Option Group

```go
group := blocks.MustOptionGroup("Category",
    option1, option2, option3,
)
```

### Confirm Dialog

```go
confirm, _ := blocks.NewConfirmDialog(
    "Are you sure?",
    blocks.MustPlainText("This action cannot be undone."),
    "Yes, delete",
    "Cancel",
    blocks.WithConfirmStyle(blocks.ConfirmStyleDanger),
)
```

### Conversation Filter

```go
filter := blocks.NewConversationFilter(
    blocks.WithFilterInclude("public", "private"),
    blocks.WithFilterExcludeExternalSharedChannels(),
    blocks.WithFilterExcludeBotUsers(),
)
```

### Dispatch Action Config

```go
config := blocks.NewDispatchActionConfig(
    blocks.DispatchOnEnterPressed,
    blocks.DispatchOnCharacterEntered,
)
```

## Error Handling

The library provides detailed validation errors:

```go
button, err := blocks.NewButton("")
if err != nil {
    // err: "text: is required: missing required field"

    var validationErr blocks.ValidationError
    if errors.As(err, &validationErr) {
        fmt.Println(validationErr.Field)   // "text"
        fmt.Println(validationErr.Message) // "is required"
    }
}
```

Common error types:
- `ErrMissingRequired` - Required field is empty
- `ErrExceedsMaxLen` - String exceeds maximum length
- `ErrExceedsMaxItems` - Array exceeds maximum items
- `ErrMinItems` - Array has fewer than minimum items

## JSON Output

All types implement `json.Marshaler`:

```go
// Single block
data, _ := json.Marshal(button)

// Builder output
jsonData, _ := builder.JSON()         // {"blocks": [...]}
jsonData, _ := builder.PrettyJSON()   // Indented
jsonData, _ := builder.BlocksJSON()   // Just the array

// Surface output
data, _ := json.Marshal(modal)
```

## Complete Example

```go
package main

import (
    "encoding/json"
    "fmt"
    "github.com/pbotsaris/goblocks/blocks"
)

func main() {
    // Create a confirmation dialog
    confirm, _ := blocks.NewConfirmDialog(
        "Confirm Submission",
        blocks.MustPlainText("Are you sure you want to submit?"),
        "Submit",
        "Cancel",
    )

    // Build a modal form
    modal := blocks.NewBuilder().
        AddHeader("User Registration").
        AddInput("Full Name", blocks.NewPlainTextInput(
            blocks.WithPlainTextInputActionID("name_input"),
            blocks.WithPlainTextInputPlaceholder("John Doe"),
        )).
        AddInput("Email", blocks.NewEmailInput(
            blocks.WithEmailInputActionID("email_input"),
        )).
        AddInput("Department", blocks.MustStaticSelect(
            []blocks.Option{
                blocks.MustOption("Engineering", "eng"),
                blocks.MustOption("Marketing", "mkt"),
                blocks.MustOption("Sales", "sales"),
            },
            blocks.WithStaticSelectActionID("dept_select"),
        )).
        AddDivider().
        AddSection(
            blocks.MustMarkdown("Please review the *terms and conditions* before submitting."),
        ).
        MustToModal("Register",
            blocks.WithModalSubmit("Submit"),
            blocks.WithModalClose("Cancel"),
            blocks.WithModalCallbackID("registration_modal"),
        )

    // Output JSON
    data, _ := json.MarshalIndent(modal, "", "  ")
    fmt.Println(string(data))
}
```

## Limits Reference

| Component | Limit |
|-----------|-------|
| Blocks per message | 50 |
| Blocks per modal/home tab | 100 |
| Actions block elements | 25 |
| Context block elements | 10 |
| Section fields | 10 |
| Select options | 100 |
| Option groups | 100 |
| Modal title | 24 characters |
| Button text | 75 characters |
| Header text | 150 characters |

---

## Socket Mode

The `socketmode` package provides a production-ready client for Slack's Socket Mode API, integrating seamlessly with the `blocks` package for type-safe responses.

### Installation

```go
import "github.com/pbotsaris/goblocks/socketmode"
```

### Quick Start

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
    // Create the client with your app-level token
    client := socketmode.New(os.Getenv("SLACK_APP_TOKEN"))

    // Register handlers
    client.OnSlashCommand(handleSlashCommand)
    client.OnInteractive(handleInteractive)
    client.OnEventsAPI(handleEvents)

    // Run with graceful shutdown
    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop()

    if err := client.Run(ctx); err != nil {
        log.Fatal(err)
    }
}

func handleSlashCommand(ctx context.Context, env *socketmode.Envelope) socketmode.Response {
    // Build a response with the blocks package
    msg := blocks.NewBuilder().
        AddSection(blocks.MustMarkdown("Hello from */mycommand*!")).
        MustToMessage("Hello!")

    return socketmode.RespondWithMessage(msg)
}

func handleInteractive(ctx context.Context, env *socketmode.Envelope) socketmode.Response {
    // Acknowledge without response
    return socketmode.NoResponse()
}

func handleEvents(ctx context.Context, env *socketmode.Envelope) socketmode.Response {
    // Process the event (env.Payload contains the raw JSON)
    return socketmode.NoResponse()
}
```

### Handler Registration

Register handlers for specific event types:

```go
// Convenience methods for common event types
client.OnSlashCommand(handler)  // slash_commands
client.OnInteractive(handler)   // interactive (buttons, modals, etc.)
client.OnEventsAPI(handler)     // events_api

// Generic handler for any event type
client.On("custom_event", handler)
```

The handler signature is:

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

### Type-Safe Responses

The `socketmode` package provides type-safe response builders that integrate with the `blocks` package:

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

### Client Options

Configure the client with functional options:

```go
client := socketmode.New(appToken,
    // Structured logging with slog
    socketmode.WithLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil))),

    // Custom metrics collection
    socketmode.WithMetrics(&myMetrics{}),

    // Limit concurrent handler execution
    socketmode.WithMaxConcurrency(20),

    // Handler execution timeout
    socketmode.WithHandlerTimeout(60 * time.Second),

    // Connection timeouts
    socketmode.WithHelloTimeout(30 * time.Second),
    socketmode.WithPingInterval(5 * time.Second),
    socketmode.WithPongTimeout(10 * time.Second),

    // Custom HTTP client
    socketmode.WithHTTPClient(&http.Client{
        Timeout: 30 * time.Second,
    }),
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

### Metrics & Observability

Implement the `MetricsHook` interface to collect metrics:

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
    handlerPanics     prometheus.Counter
    // ...
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

// ... implement other methods
```

### Running with an HTTP Server

Run Socket Mode alongside an HTTP server (e.g., for health checks or webhooks):

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

    // Socket Mode client
    client := socketmode.New(os.Getenv("SLACK_APP_TOKEN"))
    client.OnSlashCommand(handleSlashCommand)

    // HTTP server for health checks
    mux := http.NewServeMux()
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("ok"))
    })

    server := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }

    // Run both concurrently
    g, ctx := errgroup.WithContext(ctx)

    // Socket Mode
    g.Go(func() error {
        return client.Run(ctx)
    })

    // HTTP Server
    g.Go(func() error {
        log.Println("HTTP server listening on :8080")
        return server.ListenAndServe()
    })

    // Graceful shutdown
    g.Go(func() error {
        <-ctx.Done()
        log.Println("Shutting down...")
        return server.Shutdown(context.Background())
    })

    if err := g.Wait(); err != nil && err != http.ErrServerClosed {
        log.Fatal(err)
    }
}
```

### Error Handling

The client automatically classifies errors:

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

Exponential backoff is applied automatically:
- Base delay: 1 second
- Max delay: 30 seconds
- Jitter: 0-1 second
- Resets after 60 seconds of stable connection

### Panic Recovery

Handler panics are automatically recovered. The panic is logged, metrics are recorded, and an empty response is sent:

```go
client.OnSlashCommand(func(ctx context.Context, env *socketmode.Envelope) socketmode.Response {
    panic("oops") // Recovered! Connection continues.
})
```

### Complete Example: Modal Workflow

```go
package main

import (
    "context"
    "encoding/json"
    "log/slog"
    "os"
    "os/signal"
    "syscall"

    "github.com/pbotsaris/goblocks/blocks"
    "github.com/pbotsaris/goblocks/socketmode"
)

func main() {
    logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

    client := socketmode.New(
        os.Getenv("SLACK_APP_TOKEN"),
        socketmode.WithLogger(logger),
    )

    // Slash command opens a modal
    client.OnSlashCommand(func(ctx context.Context, env *socketmode.Envelope) socketmode.Response {
        // Note: Opening modals requires using the Web API with the trigger_id
        // from the payload. This example shows the response pattern.
        return socketmode.NoResponse()
    })

    // Handle modal submissions
    client.OnInteractive(func(ctx context.Context, env *socketmode.Envelope) socketmode.Response {
        // Parse the payload to determine action type
        var payload struct {
            Type string `json:"type"`
        }
        if err := json.Unmarshal(env.Payload, &payload); err != nil {
            return socketmode.NoResponse()
        }

        switch payload.Type {
        case "view_submission":
            return handleViewSubmission(env)
        case "block_actions":
            return handleBlockActions(env)
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
    // Validate form data
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

    // Example validation
    email := submission.View.State.Values["email_block"]["email_input"].Value
    if email == "" {
        return socketmode.RespondWithErrors(map[string]string{
            "email_block": "Email is required",
        })
    }

    // Success - close the modal
    return socketmode.RespondWithModalClear()
}

func handleBlockActions(env *socketmode.Envelope) socketmode.Response {
    // Handle button clicks, select changes, etc.
    return socketmode.NoResponse()
}
```

## License

MIT
