# goblocks

Go library for building Slack Block Kit UIs and communicating with Slack.

## Packages

| Package | Description | Docs |
|---------|-------------|------|
| [blocks](./blocks) | Type-safe Block Kit builder | [README](./blocks/README.md) |
| [socketmode](./socketmode) | Socket Mode client | [README](./socketmode/README.md) |

## Installation

```bash
go get github.com/pbotsaris/goblocks
```

## Quick Start

### Building Block Kit UIs

```go
import "github.com/pbotsaris/goblocks/blocks"

// Build a message with the fluent builder
message := blocks.NewBuilder().
    AddHeader("Welcome!").
    AddSection(blocks.MustMarkdown("Hello, *world*!")).
    AddDivider().
    AddActions([]blocks.ActionsElement{
        blocks.MustButton("Click me", blocks.WithActionID("btn_click")),
    }).
    MustToMessage("Welcome message")
```

### Socket Mode

```go
import "github.com/pbotsaris/goblocks/socketmode"

client := socketmode.New(os.Getenv("SLACK_APP_TOKEN"))

client.OnSlashCommand(func(ctx context.Context, env *socketmode.Envelope) socketmode.Response {
    msg := blocks.NewBuilder().
        AddSection(blocks.MustMarkdown("Hello from */mycommand*!")).
        MustToMessage("Hello!")
    return socketmode.RespondWithMessage(msg)
})

ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
defer stop()

client.Run(ctx)
```

## Roadmap

- [x] Block Kit builder (`blocks`)
- [x] Socket Mode client (`socketmode`)
- [ ] HTTP Mode (Events API via webhooks)
- [ ] Web API client (chat.postMessage, views.open, etc.)

## License

MIT
