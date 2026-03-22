package blocks

import "encoding/json"

// MarkdownBlock displays formatted markdown content.
// Available in messages only.
type MarkdownBlock struct {
	text    string
	blockID string
}

// Marker interface implementation
func (MarkdownBlock) block() {}

// MarshalJSON implements json.Marshaler.
func (m MarkdownBlock) MarshalJSON() ([]byte, error) {
	result := map[string]any{
		"type": "markdown",
		"text": m.text,
	}
	if m.blockID != "" {
		result["block_id"] = m.blockID
	}
	return json.Marshal(result)
}

// MarkdownBlockOption configures a MarkdownBlock.
type MarkdownBlockOption func(*MarkdownBlock)

// NewMarkdownBlock creates a new markdown block.
func NewMarkdownBlock(text string, opts ...MarkdownBlockOption) (MarkdownBlock, error) {
	if err := validateRequired("text", text); err != nil {
		return MarkdownBlock{}, err
	}

	m := MarkdownBlock{text: text}

	for _, opt := range opts {
		opt(&m)
	}

	return m, nil
}

// MustMarkdownBlock creates a MarkdownBlock or panics on error.
func MustMarkdownBlock(text string, opts ...MarkdownBlockOption) MarkdownBlock {
	m, err := NewMarkdownBlock(text, opts...)
	if err != nil {
		panic(err)
	}
	return m
}

// WithMarkdownBlockID sets the block_id.
func WithMarkdownBlockID(id string) MarkdownBlockOption {
	return func(m *MarkdownBlock) {
		m.blockID = id
	}
}
