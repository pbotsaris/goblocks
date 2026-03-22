package blocks

import "encoding/json"

// URLInput allows users to enter a URL into a single-line field.
type URLInput struct {
	actionID             string
	initialValue         string
	dispatchActionConfig *DispatchActionConfig
	focusOnLoad          bool
	placeholder          *PlainText
}

// Marker interface implementation
func (URLInput) inputElement() {}

// MarshalJSON implements json.Marshaler.
func (u URLInput) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type": "url_text_input",
	}
	if u.actionID != "" {
		m["action_id"] = u.actionID
	}
	if u.initialValue != "" {
		m["initial_value"] = u.initialValue
	}
	if u.dispatchActionConfig != nil {
		m["dispatch_action_config"] = u.dispatchActionConfig
	}
	if u.focusOnLoad {
		m["focus_on_load"] = true
	}
	if u.placeholder != nil {
		m["placeholder"] = u.placeholder
	}
	return json.Marshal(m)
}

// URLInputOption configures a URLInput.
type URLInputOption func(*URLInput)

// NewURLInput creates a new URL input element.
func NewURLInput(opts ...URLInputOption) URLInput {
	u := URLInput{}
	for _, opt := range opts {
		opt(&u)
	}
	return u
}

// WithURLInputActionID sets the action_id.
func WithURLInputActionID(id string) URLInputOption {
	return func(u *URLInput) {
		u.actionID = id
	}
}

// WithURLInputInitialValue sets the initial URL value.
func WithURLInputInitialValue(value string) URLInputOption {
	return func(u *URLInput) {
		u.initialValue = value
	}
}

// WithURLInputDispatchActionConfig sets when to dispatch block_actions.
func WithURLInputDispatchActionConfig(config DispatchActionConfig) URLInputOption {
	return func(u *URLInput) {
		u.dispatchActionConfig = &config
	}
}

// WithURLInputFocusOnLoad sets auto-focus.
func WithURLInputFocusOnLoad() URLInputOption {
	return func(u *URLInput) {
		u.focusOnLoad = true
	}
}

// WithURLInputPlaceholder sets placeholder text.
// Max 150 characters.
func WithURLInputPlaceholder(text string) URLInputOption {
	return func(u *URLInput) {
		if pt, err := NewPlainText(text); err == nil {
			u.placeholder = &pt
		}
	}
}
