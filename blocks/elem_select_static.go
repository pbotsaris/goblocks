package blocks

import "encoding/json"

// StaticSelect allows users to choose one option from a static list.
type StaticSelect struct {
	actionID      string
	placeholder   *PlainText
	options       []Option
	optionGroups  []OptionGroup
	initialOption *Option
	confirm       *ConfirmDialog
	focusOnLoad   bool
}

// Marker interface implementations
func (StaticSelect) sectionAccessory() {}
func (StaticSelect) actionsElement()   {}
func (StaticSelect) inputElement()     {}

// MarshalJSON implements json.Marshaler.
func (s StaticSelect) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type": "static_select",
	}
	if s.actionID != "" {
		m["action_id"] = s.actionID
	}
	if s.placeholder != nil {
		m["placeholder"] = s.placeholder
	}
	if len(s.options) > 0 {
		m["options"] = s.options
	}
	if len(s.optionGroups) > 0 {
		m["option_groups"] = s.optionGroups
	}
	if s.initialOption != nil {
		m["initial_option"] = s.initialOption
	}
	if s.confirm != nil {
		m["confirm"] = s.confirm
	}
	if s.focusOnLoad {
		m["focus_on_load"] = true
	}
	return json.Marshal(m)
}

// StaticSelectOption configures a StaticSelect.
type StaticSelectOption func(*StaticSelect)

// NewStaticSelect creates a static select with options.
// Max 100 options.
func NewStaticSelect(options []Option, opts ...StaticSelectOption) (StaticSelect, error) {
	if err := validateMinItems("options", options, 1); err != nil {
		return StaticSelect{}, err
	}
	if err := validateMaxItems("options", options, 100); err != nil {
		return StaticSelect{}, err
	}

	s := StaticSelect{options: options}

	for _, opt := range opts {
		opt(&s)
	}

	return s, nil
}

// NewStaticSelectWithGroups creates a static select with option groups.
// Max 100 groups.
func NewStaticSelectWithGroups(groups []OptionGroup, opts ...StaticSelectOption) (StaticSelect, error) {
	if err := validateMinItems("option_groups", groups, 1); err != nil {
		return StaticSelect{}, err
	}
	if err := validateMaxItems("option_groups", groups, 100); err != nil {
		return StaticSelect{}, err
	}

	s := StaticSelect{optionGroups: groups}

	for _, opt := range opts {
		opt(&s)
	}

	return s, nil
}

// WithStaticSelectActionID sets the action_id.
func WithStaticSelectActionID(id string) StaticSelectOption {
	return func(s *StaticSelect) {
		s.actionID = id
	}
}

// WithStaticSelectPlaceholder sets placeholder text.
func WithStaticSelectPlaceholder(text string) StaticSelectOption {
	return func(s *StaticSelect) {
		if pt, err := NewPlainText(text); err == nil {
			s.placeholder = &pt
		}
	}
}

// WithStaticSelectInitialOption sets the initially selected option.
func WithStaticSelectInitialOption(option Option) StaticSelectOption {
	return func(s *StaticSelect) {
		s.initialOption = &option
	}
}

// WithStaticSelectConfirm adds a confirmation dialog.
func WithStaticSelectConfirm(confirm ConfirmDialog) StaticSelectOption {
	return func(s *StaticSelect) {
		s.confirm = &confirm
	}
}

// WithStaticSelectFocusOnLoad sets auto-focus.
func WithStaticSelectFocusOnLoad() StaticSelectOption {
	return func(s *StaticSelect) {
		s.focusOnLoad = true
	}
}

// MultiStaticSelect allows users to choose multiple options from a static list.
type MultiStaticSelect struct {
	actionID       string
	placeholder    *PlainText
	options        []Option
	optionGroups   []OptionGroup
	initialOptions []Option
	confirm        *ConfirmDialog
	maxSelectedItems int
	focusOnLoad    bool
}

// Marker interface implementations
func (MultiStaticSelect) sectionAccessory() {}
func (MultiStaticSelect) actionsElement()   {}
func (MultiStaticSelect) inputElement()     {}

// MarshalJSON implements json.Marshaler.
func (m MultiStaticSelect) MarshalJSON() ([]byte, error) {
	out := map[string]any{
		"type": "multi_static_select",
	}
	if m.actionID != "" {
		out["action_id"] = m.actionID
	}
	if m.placeholder != nil {
		out["placeholder"] = m.placeholder
	}
	if len(m.options) > 0 {
		out["options"] = m.options
	}
	if len(m.optionGroups) > 0 {
		out["option_groups"] = m.optionGroups
	}
	if len(m.initialOptions) > 0 {
		out["initial_options"] = m.initialOptions
	}
	if m.confirm != nil {
		out["confirm"] = m.confirm
	}
	if m.maxSelectedItems > 0 {
		out["max_selected_items"] = m.maxSelectedItems
	}
	if m.focusOnLoad {
		out["focus_on_load"] = true
	}
	return json.Marshal(out)
}

// MultiStaticSelectOption configures a MultiStaticSelect.
type MultiStaticSelectOption func(*MultiStaticSelect)

// NewMultiStaticSelect creates a multi-select with options.
func NewMultiStaticSelect(options []Option, opts ...MultiStaticSelectOption) (MultiStaticSelect, error) {
	if err := validateMinItems("options", options, 1); err != nil {
		return MultiStaticSelect{}, err
	}
	if err := validateMaxItems("options", options, 100); err != nil {
		return MultiStaticSelect{}, err
	}

	m := MultiStaticSelect{options: options}

	for _, opt := range opts {
		opt(&m)
	}

	return m, nil
}

// NewMultiStaticSelectWithGroups creates a multi-select with option groups.
func NewMultiStaticSelectWithGroups(groups []OptionGroup, opts ...MultiStaticSelectOption) (MultiStaticSelect, error) {
	if err := validateMinItems("option_groups", groups, 1); err != nil {
		return MultiStaticSelect{}, err
	}
	if err := validateMaxItems("option_groups", groups, 100); err != nil {
		return MultiStaticSelect{}, err
	}

	m := MultiStaticSelect{optionGroups: groups}

	for _, opt := range opts {
		opt(&m)
	}

	return m, nil
}

// WithMultiStaticSelectActionID sets the action_id.
func WithMultiStaticSelectActionID(id string) MultiStaticSelectOption {
	return func(m *MultiStaticSelect) {
		m.actionID = id
	}
}

// WithMultiStaticSelectPlaceholder sets placeholder text.
func WithMultiStaticSelectPlaceholder(text string) MultiStaticSelectOption {
	return func(m *MultiStaticSelect) {
		if pt, err := NewPlainText(text); err == nil {
			m.placeholder = &pt
		}
	}
}

// WithMultiStaticSelectInitialOptions sets initially selected options.
func WithMultiStaticSelectInitialOptions(options ...Option) MultiStaticSelectOption {
	return func(m *MultiStaticSelect) {
		m.initialOptions = options
	}
}

// WithMultiStaticSelectConfirm adds a confirmation dialog.
func WithMultiStaticSelectConfirm(confirm ConfirmDialog) MultiStaticSelectOption {
	return func(m *MultiStaticSelect) {
		m.confirm = &confirm
	}
}

// WithMultiStaticSelectMaxItems sets max selectable items.
func WithMultiStaticSelectMaxItems(max int) MultiStaticSelectOption {
	return func(m *MultiStaticSelect) {
		m.maxSelectedItems = max
	}
}

// WithMultiStaticSelectFocusOnLoad sets auto-focus.
func WithMultiStaticSelectFocusOnLoad() MultiStaticSelectOption {
	return func(m *MultiStaticSelect) {
		m.focusOnLoad = true
	}
}
