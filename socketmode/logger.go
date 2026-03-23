package socketmode

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"
)

// ANSI color codes
const (
	reset     = "\033[0m"
	bold      = "\033[1m"
	dim       = "\033[2m"
	red       = "\033[31m"
	green     = "\033[32m"
	yellow    = "\033[33m"
	blue      = "\033[34m"
	magenta   = "\033[35m"
	cyan      = "\033[36m"
	white     = "\033[37m"
	boldRed   = "\033[1;31m"
	boldGreen = "\033[1;32m"
	boldYellow = "\033[1;33m"
	boldBlue  = "\033[1;34m"
	boldCyan  = "\033[1;36m"
)

// ColoredHandler is a slog.Handler that outputs colored, human-readable logs.
type ColoredHandler struct {
	opts   ColoredHandlerOptions
	mu     *sync.Mutex
	out    io.Writer
	attrs  []slog.Attr
	groups []string
}

// ColoredHandlerOptions configures the colored handler.
type ColoredHandlerOptions struct {
	// Level is the minimum level to log. Defaults to slog.LevelInfo.
	Level slog.Leveler
	// TimeFormat is the format for timestamps. Defaults to "15:04:05".
	TimeFormat string
	// ShowDate includes the date in timestamps. Defaults to false.
	ShowDate bool
}

// NewColoredHandler creates a new colored log handler.
// Writes to os.Stderr by default.
func NewColoredHandler(opts *ColoredHandlerOptions) *ColoredHandler {
	return NewColoredHandlerWithWriter(os.Stderr, opts)
}

// NewColoredHandlerWithWriter creates a new colored log handler with a custom writer.
func NewColoredHandlerWithWriter(w io.Writer, opts *ColoredHandlerOptions) *ColoredHandler {
	if opts == nil {
		opts = &ColoredHandlerOptions{}
	}
	if opts.Level == nil {
		opts.Level = slog.LevelInfo
	}
	if opts.TimeFormat == "" {
		if opts.ShowDate {
			opts.TimeFormat = "2006/01/02 15:04:05"
		} else {
			opts.TimeFormat = "15:04:05"
		}
	}

	return &ColoredHandler{
		opts: *opts,
		mu:   &sync.Mutex{},
		out:  w,
	}
}

// Enabled implements slog.Handler.
func (h *ColoredHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

// Handle implements slog.Handler.
func (h *ColoredHandler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Format timestamp
	timeStr := fmt.Sprintf("%s%s%s", dim, r.Time.Format(h.opts.TimeFormat), reset)

	// Format level with color
	levelStr := h.formatLevel(r.Level)

	// Format message
	msgStr := r.Message

	// Start building output
	fmt.Fprintf(h.out, "%s %s %s", timeStr, levelStr, msgStr)

	// Add attributes
	r.Attrs(func(a slog.Attr) bool {
		h.writeAttr(a)
		return true
	})

	// Add handler-level attrs
	for _, a := range h.attrs {
		h.writeAttr(a)
	}

	fmt.Fprintln(h.out)
	return nil
}

func (h *ColoredHandler) writeAttr(a slog.Attr) {
	if a.Equal(slog.Attr{}) {
		return
	}

	key := a.Key
	val := a.Value.Any()

	// Format duration nicely
	if d, ok := val.(time.Duration); ok {
		val = d.Round(time.Millisecond).String()
	}

	fmt.Fprintf(h.out, " %s%s%s=%v", cyan, key, reset, val)
}

func (h *ColoredHandler) formatLevel(level slog.Level) string {
	switch {
	case level >= slog.LevelError:
		return fmt.Sprintf("%sERROR%s", boldRed, reset)
	case level >= slog.LevelWarn:
		return fmt.Sprintf("%sWARN%s ", boldYellow, reset)
	case level >= slog.LevelInfo:
		return fmt.Sprintf("%sINFO%s ", boldGreen, reset)
	default:
		return fmt.Sprintf("%sDEBUG%s", dim, reset)
	}
}

// WithAttrs implements slog.Handler.
func (h *ColoredHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ColoredHandler{
		opts:   h.opts,
		mu:     h.mu,
		out:    h.out,
		attrs:  append(h.attrs, attrs...),
		groups: h.groups,
	}
}

// WithGroup implements slog.Handler.
func (h *ColoredHandler) WithGroup(name string) slog.Handler {
	return &ColoredHandler{
		opts:   h.opts,
		mu:     h.mu,
		out:    h.out,
		attrs:  h.attrs,
		groups: append(h.groups, name),
	}
}

// NewColoredLogger creates a new slog.Logger with colored output.
// This is a convenience function for quick setup.
func NewColoredLogger() *slog.Logger {
	return slog.New(NewColoredHandler(nil))
}

// NewColoredLoggerWithLevel creates a new colored logger with a specific level.
func NewColoredLoggerWithLevel(level slog.Level) *slog.Logger {
	return slog.New(NewColoredHandler(&ColoredHandlerOptions{
		Level: level,
	}))
}
