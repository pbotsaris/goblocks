package blocks

import "encoding/json"

// ConfirmStyle represents the style of the confirm button.
type ConfirmStyle string

const (
	ConfirmStylePrimary ConfirmStyle = "primary"
	ConfirmStyleDanger  ConfirmStyle = "danger"
)

// ConfirmDialog represents a confirmation dialog composition object.
// Adds a confirmation step to interactive elements.
type ConfirmDialog struct {
	title   PlainText
	text    TextObject
	confirm PlainText
	deny    PlainText
	style   ConfirmStyle
}

// MarshalJSON implements json.Marshaler.
func (c ConfirmDialog) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"title":   c.title,
		"text":    c.text,
		"confirm": c.confirm,
		"deny":    c.deny,
	}
	if c.style != "" {
		m["style"] = c.style
	}
	return json.Marshal(m)
}

// ConfirmDialogOption configures a ConfirmDialog.
type ConfirmDialogOption func(*ConfirmDialog)

// NewConfirmDialog creates a new confirmation dialog.
// title max: 100 chars, text max: 300 chars, confirm/deny max: 30 chars
func NewConfirmDialog(title, text, confirm, deny string, opts ...ConfirmDialogOption) (ConfirmDialog, error) {

	if err := validateRequiredMaxLen("title", title, 100); err != nil {
		return ConfirmDialog{}, err
	}
	if err := validateRequiredMaxLen("text", text, 300); err != nil {
		return ConfirmDialog{}, err
	}
	if err := validateRequiredMaxLen("confirm", confirm, 30); err != nil {
		return ConfirmDialog{}, err
	}
	if err := validateRequiredMaxLen("deny", deny, 30); err != nil {
		return ConfirmDialog{}, err
	}

	titlePt, err := NewPlainText(title)
	if err != nil {
		return ConfirmDialog{}, err
	}

	textPt, err := NewPlainText(text)
	if err != nil {
		return ConfirmDialog{}, err
	}

	confirmPt, err := NewPlainText(confirm)
	if err != nil {
		return ConfirmDialog{}, err
	}

	denyPt, err := NewPlainText(deny)
	if err != nil {
		return ConfirmDialog{}, err
	}

	c := ConfirmDialog{
		title:   titlePt,
		text:    textPt,
		confirm: confirmPt,
		deny:    denyPt,
	}

	for _, opt := range opts {
		opt(&c)
	}

	return c, nil
}

// NewConfirmDialogWithMarkdown creates a confirm dialog with markdown text.
func NewConfirmDialogWithMarkdown(title, text, confirm, deny string, opts ...ConfirmDialogOption) (ConfirmDialog, error) {
	if err := validateRequiredMaxLen("title", title, 100); err != nil {
		return ConfirmDialog{}, err
	}
	if err := validateRequiredMaxLen("text", text, 300); err != nil {
		return ConfirmDialog{}, err
	}
	if err := validateRequiredMaxLen("confirm", confirm, 30); err != nil {
		return ConfirmDialog{}, err
	}
	if err := validateRequiredMaxLen("deny", deny, 30); err != nil {
		return ConfirmDialog{}, err
	}

	titlePt, err := NewPlainText(title)
	if err != nil {
		return ConfirmDialog{}, err
	}

	textMd, err := NewMarkdown(text)
	if err != nil {
		return ConfirmDialog{}, err
	}

	confirmPt, err := NewPlainText(confirm)
	if err != nil {
		return ConfirmDialog{}, err
	}

	denyPt, err := NewPlainText(deny)
	if err != nil {
		return ConfirmDialog{}, err
	}

	c := ConfirmDialog{
		title:   titlePt,
		text:    textMd,
		confirm: confirmPt,
		deny:    denyPt,
	}

	for _, opt := range opts {
		opt(&c)
	}

	return c, nil
}

// WithConfirmStyle sets the style of the confirm button.
func WithConfirmStyle(style ConfirmStyle) ConfirmDialogOption {
	return func(c *ConfirmDialog) {
		c.style = style
	}
}
