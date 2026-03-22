package blocks

import "encoding/json"

// ButtonStyle represents button visual styles.
type ButtonStyle string

const (
	ButtonStylePrimary ButtonStyle = "primary"
	ButtonStyleDanger  ButtonStyle = "danger"
)

// Button represents an interactive button element.
type Button struct {
	text               PlainText
	actionID           string
	value              string
	style              ButtonStyle
	url                string
	confirm            *ConfirmDialog
	accessibilityLabel string
}

// Marker interface implementations
func (Button) sectionAccessory() {}
func (Button) actionsElement()   {}

// MarshalJSON implements json.Marshaler.
func (b Button) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type": "button",
		"text": b.text,
	}
	if b.actionID != "" {
		m["action_id"] = b.actionID
	}
	if b.value != "" {
		m["value"] = b.value
	}
	if b.style != "" {
		m["style"] = b.style
	}
	if b.url != "" {
		m["url"] = b.url
	}
	if b.confirm != nil {
		m["confirm"] = b.confirm
	}
	if b.accessibilityLabel != "" {
		m["accessibility_label"] = b.accessibilityLabel
	}
	return json.Marshal(m)
}

// ButtonOption configures a Button.
type ButtonOption func(*Button)

// NewButton creates a new Button with the given text.
// text max: 75 characters
func NewButton(text string, opts ...ButtonOption) (Button, error) {
	if err := validateRequiredMaxLen("text", text, 75); err != nil {
		return Button{}, err
	}

	pt, err := NewPlainText(text)
	if err != nil {
		return Button{}, err
	}

	b := Button{text: pt}

	for _, opt := range opts {
		opt(&b)
	}

	return b, nil
}

// MustButton creates a Button or panics on error.
func MustButton(text string, opts ...ButtonOption) Button {
	b, err := NewButton(text, opts...)
	if err != nil {
		panic(err)
	}
	return b
}

// WithActionID sets the button's action_id.
// Max 255 characters.
func WithActionID(id string) ButtonOption {
	return func(b *Button) {
		b.actionID = id
	}
}

// WithValue sets the button's value sent in interaction payload.
// Max 2000 characters.
func WithValue(value string) ButtonOption {
	return func(b *Button) {
		b.value = value
	}
}

// WithButtonStyle sets the button's visual style.
func WithButtonStyle(style ButtonStyle) ButtonOption {
	return func(b *Button) {
		b.style = style
	}
}

// WithURL sets a URL to load when the button is clicked.
// Max 3000 characters.
func WithURL(url string) ButtonOption {
	return func(b *Button) {
		b.url = url
	}
}

// WithButtonConfirm adds a confirmation dialog to the button.
func WithButtonConfirm(confirm ConfirmDialog) ButtonOption {
	return func(b *Button) {
		b.confirm = &confirm
	}
}

// WithAccessibilityLabel sets a label for screen readers.
// Max 75 characters.
func WithAccessibilityLabel(label string) ButtonOption {
	return func(b *Button) {
		b.accessibilityLabel = label
	}
}
