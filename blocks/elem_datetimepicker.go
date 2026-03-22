package blocks

import "encoding/json"

// DatetimePicker allows users to select both date and time.
type DatetimePicker struct {
	actionID        string
	initialDateTime int64
	confirm         *ConfirmDialog
	focusOnLoad     bool
}

// Marker interface implementations
// Note: DatetimePicker is only valid in input blocks
func (DatetimePicker) inputElement() {}

// MarshalJSON implements json.Marshaler.
func (d DatetimePicker) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type": "datetimepicker",
	}
	if d.actionID != "" {
		m["action_id"] = d.actionID
	}
	if d.initialDateTime != 0 {
		m["initial_date_time"] = d.initialDateTime
	}
	if d.confirm != nil {
		m["confirm"] = d.confirm
	}
	if d.focusOnLoad {
		m["focus_on_load"] = true
	}
	return json.Marshal(m)
}

// DatetimePickerOption configures a DatetimePicker.
type DatetimePickerOption func(*DatetimePicker)

// NewDatetimePicker creates a new datetime picker element.
func NewDatetimePicker(opts ...DatetimePickerOption) DatetimePicker {
	d := DatetimePicker{}

	for _, opt := range opts {
		opt(&d)
	}

	return d
}

// WithDatetimePickerActionID sets the action_id.
func WithDatetimePickerActionID(id string) DatetimePickerOption {
	return func(d *DatetimePicker) {
		d.actionID = id
	}
}

// WithInitialDateTime sets the initial date/time as Unix timestamp.
func WithInitialDateTime(timestamp int64) DatetimePickerOption {
	return func(d *DatetimePicker) {
		d.initialDateTime = timestamp
	}
}

// WithDatetimePickerConfirm adds a confirmation dialog.
func WithDatetimePickerConfirm(confirm ConfirmDialog) DatetimePickerOption {
	return func(d *DatetimePicker) {
		d.confirm = &confirm
	}
}

// WithDatetimePickerFocusOnLoad sets auto-focus.
func WithDatetimePickerFocusOnLoad() DatetimePickerOption {
	return func(d *DatetimePicker) {
		d.focusOnLoad = true
	}
}
