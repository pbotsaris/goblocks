package blocks

import "encoding/json"

// PlainText represents a plain_text composition object.
type PlainText struct {
	text  string
	emoji bool
}

// Marker interface implementations
func (PlainText) textObject()     {}
func (PlainText) plainTextOnly()  {}
func (PlainText) contextElement() {}

// MarshalJSON implements json.Marshaler.
func (p PlainText) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type  string `json:"type"`
		Text  string `json:"text"`
		Emoji bool   `json:"emoji,omitempty"`
	}{
		Type:  "plain_text",
		Text:  p.text,
		Emoji: p.emoji,
	})
}

// Text returns the text content.
func (p PlainText) Text() string {
	return p.text
}

// PlainTextOption configures a PlainText.
type PlainTextOption func(*PlainText)

// NewPlainText creates a new PlainText with the given text.
// By default, emoji rendering is enabled.
func NewPlainText(text string, opts ...PlainTextOption) (PlainText, error) {
	if err := validateRequired("text", text); err != nil {
		return PlainText{}, err
	}

	p := PlainText{
		text:  text,
		emoji: true,
	}

	for _, opt := range opts {
		opt(&p)
	}

	return p, nil
}

// MustPlainText creates a PlainText or panics on error.
func MustPlainText(text string, opts ...PlainTextOption) PlainText {
	p, err := NewPlainText(text, opts...)
	if err != nil {
		panic(err)
	}
	return p
}

// WithEmoji sets whether to render emoji codes.
func WithEmoji(emoji bool) PlainTextOption {
	return func(p *PlainText) {
		p.emoji = emoji
	}
}

// Markdown represents a mrkdwn composition object.
type Markdown struct {
	text     string
	verbatim bool
}

// Marker interface implementations
func (Markdown) textObject()     {}
func (Markdown) contextElement() {}

// MarshalJSON implements json.Marshaler.
func (m Markdown) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type     string `json:"type"`
		Text     string `json:"text"`
		Verbatim bool   `json:"verbatim,omitempty"`
	}{
		Type:     "mrkdwn",
		Text:     m.text,
		Verbatim: m.verbatim,
	})
}

// Text returns the text content.
func (m Markdown) Text() string {
	return m.text
}

// MarkdownOption configures a Markdown.
type MarkdownOption func(*Markdown)

// NewMarkdown creates a new Markdown text object.
func NewMarkdown(text string, opts ...MarkdownOption) (Markdown, error) {
	if err := validateRequired("text", text); err != nil {
		return Markdown{}, err
	}

	m := Markdown{text: text}

	for _, opt := range opts {
		opt(&m)
	}

	return m, nil
}

// MustMarkdown creates a Markdown or panics on error.
func MustMarkdown(text string, opts ...MarkdownOption) Markdown {
	m, err := NewMarkdown(text, opts...)
	if err != nil {
		panic(err)
	}
	return m
}

// WithVerbatim disables automatic parsing of URLs, mentions, and emoji.
func WithVerbatim(verbatim bool) MarkdownOption {
	return func(m *Markdown) {
		m.verbatim = verbatim
	}
}
