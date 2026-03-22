package blocks

import "encoding/json"

// RadioButtons allows users to choose one item from a list.
type RadioButtons struct {
	actionID      string
	options       []Option
	initialOption *Option
	confirm       *ConfirmDialog
	focusOnLoad   bool
}

// Marker interface implementations
func (RadioButtons) sectionAccessory() {}
func (RadioButtons) actionsElement()   {}
func (RadioButtons) inputElement()     {}

// MarshalJSON implements json.Marshaler.
func (r RadioButtons) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type":    "radio_buttons",
		"options": r.options,
	}
	if r.actionID != "" {
		m["action_id"] = r.actionID
	}
	if r.initialOption != nil {
		m["initial_option"] = r.initialOption
	}
	if r.confirm != nil {
		m["confirm"] = r.confirm
	}
	if r.focusOnLoad {
		m["focus_on_load"] = true
	}
	return json.Marshal(m)
}

// RadioButtonsOption configures a RadioButtons element.
type RadioButtonsOption func(*RadioButtons)

// NewRadioButtons creates a new radio buttons element.
// Max 10 options.
func NewRadioButtons(options []Option, opts ...RadioButtonsOption) (RadioButtons, error) {
	if err := validateMinItems("options", options, 1); err != nil {
		return RadioButtons{}, err
	}
	if err := validateMaxItems("options", options, 10); err != nil {
		return RadioButtons{}, err
	}

	r := RadioButtons{options: options}

	for _, opt := range opts {
		opt(&r)
	}

	return r, nil
}

// WithRadioButtonsActionID sets the action_id.
func WithRadioButtonsActionID(id string) RadioButtonsOption {
	return func(r *RadioButtons) {
		r.actionID = id
	}
}

// WithRadioButtonInitialOption sets the initially selected option.
func WithRadioButtonInitialOption(option Option) RadioButtonsOption {
	return func(r *RadioButtons) {
		r.initialOption = &option
	}
}

// WithRadioButtonsConfirm adds a confirmation dialog.
func WithRadioButtonsConfirm(confirm ConfirmDialog) RadioButtonsOption {
	return func(r *RadioButtons) {
		r.confirm = &confirm
	}
}

// WithRadioButtonsFocusOnLoad sets auto-focus.
func WithRadioButtonsFocusOnLoad() RadioButtonsOption {
	return func(r *RadioButtons) {
		r.focusOnLoad = true
	}
}
