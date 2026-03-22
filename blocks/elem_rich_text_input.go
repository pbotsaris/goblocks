package blocks

import "encoding/json"

// RichTextInput allows users to enter formatted text in a WYSIWYG composer.
type RichTextInput struct {
	actionID             string
	initialValue         *RichText
	dispatchActionConfig *DispatchActionConfig
	focusOnLoad          bool
	placeholder          *PlainText
}

// Marker interface implementation
func (RichTextInput) inputElement() {}

// MarshalJSON implements json.Marshaler.
func (r RichTextInput) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type": "rich_text_input",
	}
	if r.actionID != "" {
		m["action_id"] = r.actionID
	}
	if r.initialValue != nil {
		m["initial_value"] = r.initialValue
	}
	if r.dispatchActionConfig != nil {
		m["dispatch_action_config"] = r.dispatchActionConfig
	}
	if r.focusOnLoad {
		m["focus_on_load"] = true
	}
	if r.placeholder != nil {
		m["placeholder"] = r.placeholder
	}
	return json.Marshal(m)
}

// RichTextInputOption configures a RichTextInput.
type RichTextInputOption func(*RichTextInput)

// NewRichTextInput creates a new rich text input element.
func NewRichTextInput(opts ...RichTextInputOption) RichTextInput {
	r := RichTextInput{}
	for _, opt := range opts {
		opt(&r)
	}
	return r
}

// WithRichTextInputActionID sets the action_id.
func WithRichTextInputActionID(id string) RichTextInputOption {
	return func(r *RichTextInput) {
		r.actionID = id
	}
}

// WithRichTextInputInitialValue sets the initial rich text value.
func WithRichTextInputInitialValue(value *RichText) RichTextInputOption {
	return func(r *RichTextInput) {
		r.initialValue = value
	}
}

// WithRichTextInputDispatchActionConfig sets when to dispatch block_actions.
func WithRichTextInputDispatchActionConfig(config DispatchActionConfig) RichTextInputOption {
	return func(r *RichTextInput) {
		r.dispatchActionConfig = &config
	}
}

// WithRichTextInputFocusOnLoad sets auto-focus.
func WithRichTextInputFocusOnLoad() RichTextInputOption {
	return func(r *RichTextInput) {
		r.focusOnLoad = true
	}
}

// WithRichTextInputPlaceholder sets placeholder text.
// Max 150 characters.
func WithRichTextInputPlaceholder(text string) RichTextInputOption {
	return func(r *RichTextInput) {
		if pt, err := NewPlainText(text); err == nil {
			r.placeholder = &pt
		}
	}
}
