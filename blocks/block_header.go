package blocks

import "encoding/json"

// Header displays larger-sized text.
type Header struct {
	text    PlainText
	blockID string
}

// Marker interface implementation
func (Header) block() {}

// MarshalJSON implements json.Marshaler.
func (h Header) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type": "header",
		"text": h.text,
	}
	if h.blockID != "" {
		m["block_id"] = h.blockID
	}
	return json.Marshal(m)
}

// HeaderOption configures a Header.
type HeaderOption func(*Header)

// NewHeader creates a new header block.
// text max: 150 characters, must be plain_text
func NewHeader(text string, opts ...HeaderOption) (Header, error) {
	if err := validateRequiredMaxLen("text", text, 150); err != nil {
		return Header{}, err
	}

	pt, err := NewPlainText(text)
	if err != nil {
		return Header{}, err
	}

	h := Header{text: pt}

	for _, opt := range opts {
		opt(&h)
	}

	return h, nil
}

// MustHeader creates a Header or panics on error.
func MustHeader(text string, opts ...HeaderOption) Header {
	h, err := NewHeader(text, opts...)
	if err != nil {
		panic(err)
	}
	return h
}

// WithHeaderBlockID sets the block_id.
func WithHeaderBlockID(id string) HeaderOption {
	return func(h *Header) {
		h.blockID = id
	}
}
