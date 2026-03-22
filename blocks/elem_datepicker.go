package blocks

import "encoding/json"

// DatePicker allows users to select a date from a calendar UI.
type DatePicker struct {
	actionID    string
	initialDate string
	placeholder *PlainText
	confirm     *ConfirmDialog
	focusOnLoad bool
}

// Marker interface implementations
func (DatePicker) sectionAccessory() {}
func (DatePicker) actionsElement()   {}
func (DatePicker) inputElement()     {}

// MarshalJSON implements json.Marshaler.
func (d DatePicker) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type": "datepicker",
	}
	if d.actionID != "" {
		m["action_id"] = d.actionID
	}
	if d.initialDate != "" {
		m["initial_date"] = d.initialDate
	}
	if d.placeholder != nil {
		m["placeholder"] = d.placeholder
	}
	if d.confirm != nil {
		m["confirm"] = d.confirm
	}
	if d.focusOnLoad {
		m["focus_on_load"] = true
	}
	return json.Marshal(m)
}

// DatePickerOption configures a DatePicker.
type DatePickerOption func(*DatePicker)

// NewDatePicker creates a new date picker element.
func NewDatePicker(opts ...DatePickerOption) DatePicker {
	d := DatePicker{}

	for _, opt := range opts {
		opt(&d)
	}

	return d
}

// WithDatePickerActionID sets the action_id.
func WithDatePickerActionID(id string) DatePickerOption {
	return func(d *DatePicker) {
		d.actionID = id
	}
}

// WithInitialDate sets the initially selected date.
// Format: YYYY-MM-DD
func WithInitialDate(date string) DatePickerOption {
	return func(d *DatePicker) {
		d.initialDate = date
	}
}

// WithDatePickerPlaceholder sets placeholder text.
// Max 150 characters.
func WithDatePickerPlaceholder(text string) DatePickerOption {
	return func(d *DatePicker) {
		if pt, err := NewPlainText(text); err == nil {
			d.placeholder = &pt
		}
	}
}

// WithDatePickerConfirm adds a confirmation dialog.
func WithDatePickerConfirm(confirm ConfirmDialog) DatePickerOption {
	return func(d *DatePicker) {
		d.confirm = &confirm
	}
}

// WithDatePickerFocusOnLoad sets auto-focus.
func WithDatePickerFocusOnLoad() DatePickerOption {
	return func(d *DatePicker) {
		d.focusOnLoad = true
	}
}
