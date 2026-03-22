package blocks

import "encoding/json"

// PlainTextInput allows users to enter freeform text.
type PlainTextInput struct {
	actionID             string
	placeholder          *PlainText
	initialValue         string
	multiline            bool
	minLength            int
	maxLength            int
	dispatchActionConfig *DispatchActionConfig
	focusOnLoad          bool
}

// Marker interface implementations
func (PlainTextInput) inputElement() {}

// MarshalJSON implements json.Marshaler.
func (p PlainTextInput) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type": "plain_text_input",
	}
	if p.actionID != "" {
		m["action_id"] = p.actionID
	}
	if p.placeholder != nil {
		m["placeholder"] = p.placeholder
	}
	if p.initialValue != "" {
		m["initial_value"] = p.initialValue
	}
	if p.multiline {
		m["multiline"] = true
	}
	if p.minLength > 0 {
		m["min_length"] = p.minLength
	}
	if p.maxLength > 0 {
		m["max_length"] = p.maxLength
	}
	if p.dispatchActionConfig != nil {
		m["dispatch_action_config"] = p.dispatchActionConfig
	}
	if p.focusOnLoad {
		m["focus_on_load"] = true
	}
	return json.Marshal(m)
}

// PlainTextInputOption configures a PlainTextInput.
type PlainTextInputOption func(*PlainTextInput)

// NewPlainTextInput creates a new plain text input element.
func NewPlainTextInput(opts ...PlainTextInputOption) PlainTextInput {
	p := PlainTextInput{}

	for _, opt := range opts {
		opt(&p)
	}

	return p
}

// WithPlainTextInputActionID sets the action_id.
func WithPlainTextInputActionID(id string) PlainTextInputOption {
	return func(p *PlainTextInput) {
		p.actionID = id
	}
}

// WithPlainTextInputPlaceholder sets placeholder text.
// Max 150 characters.
func WithPlainTextInputPlaceholder(text string) PlainTextInputOption {
	return func(p *PlainTextInput) {
		if pt, err := NewPlainText(text); err == nil {
			p.placeholder = &pt
		}
	}
}

// WithInitialValue sets the initial text value.
func WithInitialValue(value string) PlainTextInputOption {
	return func(p *PlainTextInput) {
		p.initialValue = value
	}
}

// WithMultiline enables multiline input.
func WithMultiline() PlainTextInputOption {
	return func(p *PlainTextInput) {
		p.multiline = true
	}
}

// WithMinLength sets minimum input length.
func WithMinLength(min int) PlainTextInputOption {
	return func(p *PlainTextInput) {
		p.minLength = min
	}
}

// WithMaxLength sets maximum input length.
func WithMaxLength(max int) PlainTextInputOption {
	return func(p *PlainTextInput) {
		p.maxLength = max
	}
}

// WithDispatchActionConfig sets when to dispatch block_actions.
func WithDispatchActionConfig(config DispatchActionConfig) PlainTextInputOption {
	return func(p *PlainTextInput) {
		p.dispatchActionConfig = &config
	}
}

// WithPlainTextInputFocusOnLoad sets auto-focus.
func WithPlainTextInputFocusOnLoad() PlainTextInputOption {
	return func(p *PlainTextInput) {
		p.focusOnLoad = true
	}
}
