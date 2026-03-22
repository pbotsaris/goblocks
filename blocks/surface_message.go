package blocks

import "encoding/json"

// Message represents a message surface.
type Message struct {
	blocks   []Block
	text     string // fallback text
	threadTS string
	mrkdwn   bool
}

// MarshalJSON implements json.Marshaler.
func (m Message) MarshalJSON() ([]byte, error) {
	out := map[string]any{
		"blocks": m.blocks,
	}
	if m.text != "" {
		out["text"] = m.text
	}
	if m.threadTS != "" {
		out["thread_ts"] = m.threadTS
	}
	if m.mrkdwn {
		out["mrkdwn"] = true
	}
	return json.Marshal(out)
}

// Blocks returns the message blocks.
func (m Message) Blocks() []Block {
	return m.blocks
}

// MessageOption configures a Message.
type MessageOption func(*Message)

// NewMessage creates a new message.
// text is fallback for notifications, max 50 blocks
func NewMessage(text string, blocks []Block, opts ...MessageOption) (Message, error) {
	if err := validateMaxItems("blocks", blocks, 50); err != nil {
		return Message{}, err
	}

	m := Message{
		text:   text,
		blocks: blocks,
	}

	for _, opt := range opts {
		opt(&m)
	}

	return m, nil
}

// NewMessageWithBlocks creates a message with blocks only.
func NewMessageWithBlocks(blocks []Block, opts ...MessageOption) (Message, error) {
	return NewMessage("", blocks, opts...)
}

// MustMessage creates a Message or panics on error.
func MustMessage(text string, blocks []Block, opts ...MessageOption) Message {
	m, err := NewMessage(text, blocks, opts...)
	if err != nil {
		panic(err)
	}
	return m
}

// WithMessageThreadTS sets the thread timestamp for replies.
func WithMessageThreadTS(ts string) MessageOption {
	return func(m *Message) {
		m.threadTS = ts
	}
}

// WithMessageMrkdwn enables markdown parsing for the text.
func WithMessageMrkdwn() MessageOption {
	return func(m *Message) {
		m.mrkdwn = true
	}
}
