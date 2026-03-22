package blocks

import "encoding/json"

// Checkboxes allows users to choose multiple items from a list.
type Checkboxes struct {
	actionID       string
	options        []Option
	initialOptions []Option
	confirm        *ConfirmDialog
	focusOnLoad    bool
}

// Marker interface implementations
func (Checkboxes) sectionAccessory() {}
func (Checkboxes) actionsElement()   {}
func (Checkboxes) inputElement()     {}

// MarshalJSON implements json.Marshaler.
func (c Checkboxes) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type":    "checkboxes",
		"options": c.options,
	}
	if c.actionID != "" {
		m["action_id"] = c.actionID
	}
	if len(c.initialOptions) > 0 {
		m["initial_options"] = c.initialOptions
	}
	if c.confirm != nil {
		m["confirm"] = c.confirm
	}
	if c.focusOnLoad {
		m["focus_on_load"] = true
	}
	return json.Marshal(m)
}

// CheckboxesOption configures a Checkboxes element.
type CheckboxesOption func(*Checkboxes)

// NewCheckboxes creates a new checkboxes element.
// Max 10 options.
func NewCheckboxes(options []Option, opts ...CheckboxesOption) (Checkboxes, error) {
	if err := validateMinItems("options", options, 1); err != nil {
		return Checkboxes{}, err
	}
	if err := validateMaxItems("options", options, 10); err != nil {
		return Checkboxes{}, err
	}

	c := Checkboxes{options: options}

	for _, opt := range opts {
		opt(&c)
	}

	return c, nil
}

// WithCheckboxesActionID sets the action_id.
func WithCheckboxesActionID(id string) CheckboxesOption {
	return func(c *Checkboxes) {
		c.actionID = id
	}
}

// WithInitialOptions sets initially selected options.
func WithInitialOptions(options ...Option) CheckboxesOption {
	return func(c *Checkboxes) {
		c.initialOptions = options
	}
}

// WithCheckboxesConfirm adds a confirmation dialog.
func WithCheckboxesConfirm(confirm ConfirmDialog) CheckboxesOption {
	return func(c *Checkboxes) {
		c.confirm = &confirm
	}
}

// WithCheckboxesFocusOnLoad sets auto-focus.
func WithCheckboxesFocusOnLoad() CheckboxesOption {
	return func(c *Checkboxes) {
		c.focusOnLoad = true
	}
}
