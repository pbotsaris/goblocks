package blocks

import "encoding/json"

// Overflow represents an overflow menu element.
type Overflow struct {
	actionID string
	options  []Option
	confirm  *ConfirmDialog
}

// Marker interface implementations
func (Overflow) sectionAccessory() {}
func (Overflow) actionsElement()   {}

// MarshalJSON implements json.Marshaler.
func (o Overflow) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type":    "overflow",
		"options": o.options,
	}
	if o.actionID != "" {
		m["action_id"] = o.actionID
	}
	if o.confirm != nil {
		m["confirm"] = o.confirm
	}
	return json.Marshal(m)
}

// OverflowOption configures an Overflow.
type OverflowOption func(*Overflow)

// NewOverflow creates a new overflow menu.
// Requires 1-5 options.
func NewOverflow(options []Option, opts ...OverflowOption) (Overflow, error) {
	if err := validateMinItems("options", options, 1); err != nil {
		return Overflow{}, err
	}
	if err := validateMaxItems("options", options, 5); err != nil {
		return Overflow{}, err
	}

	o := Overflow{options: options}

	for _, opt := range opts {
		opt(&o)
	}

	return o, nil
}

// WithOverflowActionID sets the action_id.
func WithOverflowActionID(id string) OverflowOption {
	return func(o *Overflow) {
		o.actionID = id
	}
}

// WithOverflowConfirm adds a confirmation dialog.
func WithOverflowConfirm(confirm ConfirmDialog) OverflowOption {
	return func(o *Overflow) {
		o.confirm = &confirm
	}
}
