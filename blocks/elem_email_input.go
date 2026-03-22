package blocks

import "encoding/json"

// EmailInput allows users to enter an email into a single-line field.
type EmailInput struct {
	actionID             string
	initialValue         string
	dispatchActionConfig *DispatchActionConfig
	focusOnLoad          bool
	placeholder          *PlainText
}

// Marker interface implementation
func (EmailInput) inputElement() {}

// MarshalJSON implements json.Marshaler.
func (e EmailInput) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type": "email_text_input",
	}
	if e.actionID != "" {
		m["action_id"] = e.actionID
	}
	if e.initialValue != "" {
		m["initial_value"] = e.initialValue
	}
	if e.dispatchActionConfig != nil {
		m["dispatch_action_config"] = e.dispatchActionConfig
	}
	if e.focusOnLoad {
		m["focus_on_load"] = true
	}
	if e.placeholder != nil {
		m["placeholder"] = e.placeholder
	}
	return json.Marshal(m)
}

// EmailInputOption configures an EmailInput.
type EmailInputOption func(*EmailInput)

// NewEmailInput creates a new email input element.
func NewEmailInput(opts ...EmailInputOption) EmailInput {
	e := EmailInput{}
	for _, opt := range opts {
		opt(&e)
	}
	return e
}

// WithEmailInputActionID sets the action_id.
func WithEmailInputActionID(id string) EmailInputOption {
	return func(e *EmailInput) {
		e.actionID = id
	}
}

// WithEmailInputInitialValue sets the initial email value.
func WithEmailInputInitialValue(value string) EmailInputOption {
	return func(e *EmailInput) {
		e.initialValue = value
	}
}

// WithEmailInputDispatchActionConfig sets when to dispatch block_actions.
func WithEmailInputDispatchActionConfig(config DispatchActionConfig) EmailInputOption {
	return func(e *EmailInput) {
		e.dispatchActionConfig = &config
	}
}

// WithEmailInputFocusOnLoad sets auto-focus.
func WithEmailInputFocusOnLoad() EmailInputOption {
	return func(e *EmailInput) {
		e.focusOnLoad = true
	}
}

// WithEmailInputPlaceholder sets placeholder text.
// Max 150 characters.
func WithEmailInputPlaceholder(text string) EmailInputOption {
	return func(e *EmailInput) {
		if pt, err := NewPlainText(text); err == nil {
			e.placeholder = &pt
		}
	}
}
