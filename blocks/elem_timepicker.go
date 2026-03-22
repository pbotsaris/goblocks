package blocks

import "encoding/json"

// TimePicker allows users to select a time.
type TimePicker struct {
	actionID    string
	initialTime string
	placeholder *PlainText
	confirm     *ConfirmDialog
	focusOnLoad bool
	timezone    string
}

// Marker interface implementations
func (TimePicker) sectionAccessory() {}
func (TimePicker) actionsElement()   {}
func (TimePicker) inputElement()     {}

// MarshalJSON implements json.Marshaler.
func (t TimePicker) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type": "timepicker",
	}
	if t.actionID != "" {
		m["action_id"] = t.actionID
	}
	if t.initialTime != "" {
		m["initial_time"] = t.initialTime
	}
	if t.placeholder != nil {
		m["placeholder"] = t.placeholder
	}
	if t.confirm != nil {
		m["confirm"] = t.confirm
	}
	if t.focusOnLoad {
		m["focus_on_load"] = true
	}
	if t.timezone != "" {
		m["timezone"] = t.timezone
	}
	return json.Marshal(m)
}

// TimePickerOption configures a TimePicker.
type TimePickerOption func(*TimePicker)

// NewTimePicker creates a new time picker element.
func NewTimePicker(opts ...TimePickerOption) TimePicker {
	t := TimePicker{}

	for _, opt := range opts {
		opt(&t)
	}

	return t
}

// WithTimePickerActionID sets the action_id.
func WithTimePickerActionID(id string) TimePickerOption {
	return func(t *TimePicker) {
		t.actionID = id
	}
}

// WithInitialTime sets the initially selected time.
// Format: HH:mm
func WithInitialTime(time string) TimePickerOption {
	return func(t *TimePicker) {
		t.initialTime = time
	}
}

// WithTimePickerPlaceholder sets placeholder text.
func WithTimePickerPlaceholder(text string) TimePickerOption {
	return func(t *TimePicker) {
		if pt, err := NewPlainText(text); err == nil {
			t.placeholder = &pt
		}
	}
}

// WithTimePickerConfirm adds a confirmation dialog.
func WithTimePickerConfirm(confirm ConfirmDialog) TimePickerOption {
	return func(t *TimePicker) {
		t.confirm = &confirm
	}
}

// WithTimePickerFocusOnLoad sets auto-focus.
func WithTimePickerFocusOnLoad() TimePickerOption {
	return func(t *TimePicker) {
		t.focusOnLoad = true
	}
}

// WithTimezone sets the timezone for the time picker.
func WithTimezone(tz string) TimePickerOption {
	return func(t *TimePicker) {
		t.timezone = tz
	}
}
