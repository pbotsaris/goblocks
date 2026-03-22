package blocks

import "encoding/json"

// NumberInput allows users to enter a number into a single-line field.
type NumberInput struct {
	actionID             string
	isDecimalAllowed     bool
	initialValue         string
	minValue             string
	maxValue             string
	dispatchActionConfig *DispatchActionConfig
	focusOnLoad          bool
	placeholder          *PlainText
}

// Marker interface implementation
func (NumberInput) inputElement() {}

// MarshalJSON implements json.Marshaler.
func (n NumberInput) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type":               "number_input",
		"is_decimal_allowed": n.isDecimalAllowed,
	}
	if n.actionID != "" {
		m["action_id"] = n.actionID
	}
	if n.initialValue != "" {
		m["initial_value"] = n.initialValue
	}
	if n.minValue != "" {
		m["min_value"] = n.minValue
	}
	if n.maxValue != "" {
		m["max_value"] = n.maxValue
	}
	if n.dispatchActionConfig != nil {
		m["dispatch_action_config"] = n.dispatchActionConfig
	}
	if n.focusOnLoad {
		m["focus_on_load"] = true
	}
	if n.placeholder != nil {
		m["placeholder"] = n.placeholder
	}
	return json.Marshal(m)
}

// NumberInputOption configures a NumberInput.
type NumberInputOption func(*NumberInput)

// NewNumberInput creates a new number input element.
// isDecimalAllowed specifies whether decimal numbers are allowed.
func NewNumberInput(isDecimalAllowed bool, opts ...NumberInputOption) NumberInput {
	n := NumberInput{isDecimalAllowed: isDecimalAllowed}
	for _, opt := range opts {
		opt(&n)
	}
	return n
}

// WithNumberInputActionID sets the action_id.
func WithNumberInputActionID(id string) NumberInputOption {
	return func(n *NumberInput) {
		n.actionID = id
	}
}

// WithNumberInputInitialValue sets the initial number value.
func WithNumberInputInitialValue(value string) NumberInputOption {
	return func(n *NumberInput) {
		n.initialValue = value
	}
}

// WithNumberInputMinValue sets the minimum value.
func WithNumberInputMinValue(value string) NumberInputOption {
	return func(n *NumberInput) {
		n.minValue = value
	}
}

// WithNumberInputMaxValue sets the maximum value.
func WithNumberInputMaxValue(value string) NumberInputOption {
	return func(n *NumberInput) {
		n.maxValue = value
	}
}

// WithNumberInputDispatchActionConfig sets when to dispatch block_actions.
func WithNumberInputDispatchActionConfig(config DispatchActionConfig) NumberInputOption {
	return func(n *NumberInput) {
		n.dispatchActionConfig = &config
	}
}

// WithNumberInputFocusOnLoad sets auto-focus.
func WithNumberInputFocusOnLoad() NumberInputOption {
	return func(n *NumberInput) {
		n.focusOnLoad = true
	}
}

// WithNumberInputPlaceholder sets placeholder text.
// Max 150 characters.
func WithNumberInputPlaceholder(text string) NumberInputOption {
	return func(n *NumberInput) {
		if pt, err := NewPlainText(text); err == nil {
			n.placeholder = &pt
		}
	}
}
