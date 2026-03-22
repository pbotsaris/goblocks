package blocks

import "encoding/json"

// ExternalSelect allows users to choose from options loaded from an external source.
type ExternalSelect struct {
	actionID       string
	placeholder    *PlainText
	initialOption  *Option
	minQueryLength int
	confirm        *ConfirmDialog
	focusOnLoad    bool
}

// Marker interface implementations
func (ExternalSelect) sectionAccessory() {}
func (ExternalSelect) actionsElement()   {}
func (ExternalSelect) inputElement()     {}

// MarshalJSON implements json.Marshaler.
func (e ExternalSelect) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type": "external_select",
	}
	if e.actionID != "" {
		m["action_id"] = e.actionID
	}
	if e.placeholder != nil {
		m["placeholder"] = e.placeholder
	}
	if e.initialOption != nil {
		m["initial_option"] = e.initialOption
	}
	if e.minQueryLength > 0 {
		m["min_query_length"] = e.minQueryLength
	}
	if e.confirm != nil {
		m["confirm"] = e.confirm
	}
	if e.focusOnLoad {
		m["focus_on_load"] = true
	}
	return json.Marshal(m)
}

// ExternalSelectOption configures an ExternalSelect.
type ExternalSelectOption func(*ExternalSelect)

// NewExternalSelect creates an external data source select.
func NewExternalSelect(opts ...ExternalSelectOption) ExternalSelect {
	e := ExternalSelect{}

	for _, opt := range opts {
		opt(&e)
	}

	return e
}

// WithExternalSelectActionID sets the action_id.
func WithExternalSelectActionID(id string) ExternalSelectOption {
	return func(e *ExternalSelect) {
		e.actionID = id
	}
}

// WithExternalSelectPlaceholder sets placeholder text.
func WithExternalSelectPlaceholder(text string) ExternalSelectOption {
	return func(e *ExternalSelect) {
		if pt, err := NewPlainText(text); err == nil {
			e.placeholder = &pt
		}
	}
}

// WithExternalSelectInitialOption sets the initially selected option.
func WithExternalSelectInitialOption(option Option) ExternalSelectOption {
	return func(e *ExternalSelect) {
		e.initialOption = &option
	}
}

// WithMinQueryLength sets minimum characters before querying.
func WithMinQueryLength(length int) ExternalSelectOption {
	return func(e *ExternalSelect) {
		e.minQueryLength = length
	}
}

// WithExternalSelectConfirm adds a confirmation dialog.
func WithExternalSelectConfirm(confirm ConfirmDialog) ExternalSelectOption {
	return func(e *ExternalSelect) {
		e.confirm = &confirm
	}
}

// WithExternalSelectFocusOnLoad sets auto-focus.
func WithExternalSelectFocusOnLoad() ExternalSelectOption {
	return func(e *ExternalSelect) {
		e.focusOnLoad = true
	}
}

// MultiExternalSelect allows multiple selections from an external source.
type MultiExternalSelect struct {
	actionID         string
	placeholder      *PlainText
	initialOptions   []Option
	minQueryLength   int
	maxSelectedItems int
	confirm          *ConfirmDialog
	focusOnLoad      bool
}

// Marker interface implementations
func (MultiExternalSelect) sectionAccessory() {}
func (MultiExternalSelect) actionsElement()   {}
func (MultiExternalSelect) inputElement()     {}

// MarshalJSON implements json.Marshaler.
func (m MultiExternalSelect) MarshalJSON() ([]byte, error) {
	out := map[string]any{
		"type": "multi_external_select",
	}
	if m.actionID != "" {
		out["action_id"] = m.actionID
	}
	if m.placeholder != nil {
		out["placeholder"] = m.placeholder
	}
	if len(m.initialOptions) > 0 {
		out["initial_options"] = m.initialOptions
	}
	if m.minQueryLength > 0 {
		out["min_query_length"] = m.minQueryLength
	}
	if m.maxSelectedItems > 0 {
		out["max_selected_items"] = m.maxSelectedItems
	}
	if m.confirm != nil {
		out["confirm"] = m.confirm
	}
	if m.focusOnLoad {
		out["focus_on_load"] = true
	}
	return json.Marshal(out)
}

// MultiExternalSelectOption configures a MultiExternalSelect.
type MultiExternalSelectOption func(*MultiExternalSelect)

// NewMultiExternalSelect creates a multi-select with external data source.
func NewMultiExternalSelect(opts ...MultiExternalSelectOption) MultiExternalSelect {
	m := MultiExternalSelect{}

	for _, opt := range opts {
		opt(&m)
	}

	return m
}

// WithMultiExternalSelectActionID sets the action_id.
func WithMultiExternalSelectActionID(id string) MultiExternalSelectOption {
	return func(m *MultiExternalSelect) {
		m.actionID = id
	}
}

// WithMultiExternalSelectPlaceholder sets placeholder text.
func WithMultiExternalSelectPlaceholder(text string) MultiExternalSelectOption {
	return func(m *MultiExternalSelect) {
		if pt, err := NewPlainText(text); err == nil {
			m.placeholder = &pt
		}
	}
}

// WithMultiExternalSelectInitialOptions sets initially selected options.
func WithMultiExternalSelectInitialOptions(options ...Option) MultiExternalSelectOption {
	return func(m *MultiExternalSelect) {
		m.initialOptions = options
	}
}

// WithMultiExternalSelectMinQueryLength sets minimum characters before querying.
func WithMultiExternalSelectMinQueryLength(length int) MultiExternalSelectOption {
	return func(m *MultiExternalSelect) {
		m.minQueryLength = length
	}
}

// WithMultiExternalSelectMaxItems sets max selectable items.
func WithMultiExternalSelectMaxItems(max int) MultiExternalSelectOption {
	return func(m *MultiExternalSelect) {
		m.maxSelectedItems = max
	}
}

// WithMultiExternalSelectConfirm adds a confirmation dialog.
func WithMultiExternalSelectConfirm(confirm ConfirmDialog) MultiExternalSelectOption {
	return func(m *MultiExternalSelect) {
		m.confirm = &confirm
	}
}

// WithMultiExternalSelectFocusOnLoad sets auto-focus.
func WithMultiExternalSelectFocusOnLoad() MultiExternalSelectOption {
	return func(m *MultiExternalSelect) {
		m.focusOnLoad = true
	}
}
