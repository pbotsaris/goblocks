package blocks

import "encoding/json"

// ChannelsSelect allows users to select a public channel.
type ChannelsSelect struct {
	actionID           string
	placeholder        *PlainText
	initialChannel     string
	confirm            *ConfirmDialog
	focusOnLoad        bool
	responseURLEnabled bool
}

// Marker interface implementations
func (ChannelsSelect) sectionAccessory() {}
func (ChannelsSelect) actionsElement()   {}
func (ChannelsSelect) inputElement()     {}

// MarshalJSON implements json.Marshaler.
func (c ChannelsSelect) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type": "channels_select",
	}
	if c.actionID != "" {
		m["action_id"] = c.actionID
	}
	if c.placeholder != nil {
		m["placeholder"] = c.placeholder
	}
	if c.initialChannel != "" {
		m["initial_channel"] = c.initialChannel
	}
	if c.confirm != nil {
		m["confirm"] = c.confirm
	}
	if c.focusOnLoad {
		m["focus_on_load"] = true
	}
	if c.responseURLEnabled {
		m["response_url_enabled"] = true
	}
	return json.Marshal(m)
}

// ChannelsSelectOption configures a ChannelsSelect.
type ChannelsSelectOption func(*ChannelsSelect)

// NewChannelsSelect creates a channel select element.
func NewChannelsSelect(opts ...ChannelsSelectOption) ChannelsSelect {
	c := ChannelsSelect{}

	for _, opt := range opts {
		opt(&c)
	}

	return c
}

// WithChannelsSelectActionID sets the action_id.
func WithChannelsSelectActionID(id string) ChannelsSelectOption {
	return func(c *ChannelsSelect) {
		c.actionID = id
	}
}

// WithChannelsSelectPlaceholder sets placeholder text.
func WithChannelsSelectPlaceholder(text string) ChannelsSelectOption {
	return func(c *ChannelsSelect) {
		if pt, err := NewPlainText(text); err == nil {
			c.placeholder = &pt
		}
	}
}

// WithInitialChannel sets the initially selected channel ID.
func WithInitialChannel(channelID string) ChannelsSelectOption {
	return func(c *ChannelsSelect) {
		c.initialChannel = channelID
	}
}

// WithChannelsSelectConfirm adds a confirmation dialog.
func WithChannelsSelectConfirm(confirm ConfirmDialog) ChannelsSelectOption {
	return func(c *ChannelsSelect) {
		c.confirm = &confirm
	}
}

// WithChannelsSelectFocusOnLoad sets auto-focus.
func WithChannelsSelectFocusOnLoad() ChannelsSelectOption {
	return func(c *ChannelsSelect) {
		c.focusOnLoad = true
	}
}

// WithChannelsSelectResponseURLEnabled enables response URL for workflows.
func WithChannelsSelectResponseURLEnabled() ChannelsSelectOption {
	return func(c *ChannelsSelect) {
		c.responseURLEnabled = true
	}
}

// MultiChannelsSelect allows selecting multiple public channels.
type MultiChannelsSelect struct {
	actionID         string
	placeholder      *PlainText
	initialChannels  []string
	maxSelectedItems int
	confirm          *ConfirmDialog
	focusOnLoad      bool
}

// Marker interface implementations
func (MultiChannelsSelect) sectionAccessory() {}
func (MultiChannelsSelect) actionsElement()   {}
func (MultiChannelsSelect) inputElement()     {}

// MarshalJSON implements json.Marshaler.
func (m MultiChannelsSelect) MarshalJSON() ([]byte, error) {
	out := map[string]any{
		"type": "multi_channels_select",
	}
	if m.actionID != "" {
		out["action_id"] = m.actionID
	}
	if m.placeholder != nil {
		out["placeholder"] = m.placeholder
	}
	if len(m.initialChannels) > 0 {
		out["initial_channels"] = m.initialChannels
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

// MultiChannelsSelectOption configures a MultiChannelsSelect.
type MultiChannelsSelectOption func(*MultiChannelsSelect)

// NewMultiChannelsSelect creates a multi-channel select element.
func NewMultiChannelsSelect(opts ...MultiChannelsSelectOption) MultiChannelsSelect {
	m := MultiChannelsSelect{}

	for _, opt := range opts {
		opt(&m)
	}

	return m
}

// WithMultiChannelsSelectActionID sets the action_id.
func WithMultiChannelsSelectActionID(id string) MultiChannelsSelectOption {
	return func(m *MultiChannelsSelect) {
		m.actionID = id
	}
}

// WithMultiChannelsSelectPlaceholder sets placeholder text.
func WithMultiChannelsSelectPlaceholder(text string) MultiChannelsSelectOption {
	return func(m *MultiChannelsSelect) {
		if pt, err := NewPlainText(text); err == nil {
			m.placeholder = &pt
		}
	}
}

// WithInitialChannels sets initially selected channel IDs.
func WithInitialChannels(channelIDs ...string) MultiChannelsSelectOption {
	return func(m *MultiChannelsSelect) {
		m.initialChannels = channelIDs
	}
}

// WithMultiChannelsSelectMaxItems sets max selectable items.
func WithMultiChannelsSelectMaxItems(max int) MultiChannelsSelectOption {
	return func(m *MultiChannelsSelect) {
		m.maxSelectedItems = max
	}
}

// WithMultiChannelsSelectConfirm adds a confirmation dialog.
func WithMultiChannelsSelectConfirm(confirm ConfirmDialog) MultiChannelsSelectOption {
	return func(m *MultiChannelsSelect) {
		m.confirm = &confirm
	}
}

// WithMultiChannelsSelectFocusOnLoad sets auto-focus.
func WithMultiChannelsSelectFocusOnLoad() MultiChannelsSelectOption {
	return func(m *MultiChannelsSelect) {
		m.focusOnLoad = true
	}
}
