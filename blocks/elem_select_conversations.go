package blocks

import "encoding/json"

// ConversationsSelect allows users to select a conversation.
type ConversationsSelect struct {
	actionID                   string
	placeholder                *PlainText
	initialConversation        string
	defaultToCurrentConversation bool
	filter                     *ConversationFilter
	confirm                    *ConfirmDialog
	focusOnLoad                bool
	responseURLEnabled         bool
}

// Marker interface implementations
func (ConversationsSelect) sectionAccessory() {}
func (ConversationsSelect) actionsElement()   {}
func (ConversationsSelect) inputElement()     {}

// MarshalJSON implements json.Marshaler.
func (c ConversationsSelect) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type": "conversations_select",
	}
	if c.actionID != "" {
		m["action_id"] = c.actionID
	}
	if c.placeholder != nil {
		m["placeholder"] = c.placeholder
	}
	if c.initialConversation != "" {
		m["initial_conversation"] = c.initialConversation
	}
	if c.defaultToCurrentConversation {
		m["default_to_current_conversation"] = true
	}
	if c.filter != nil {
		m["filter"] = c.filter
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

// ConversationsSelectOption configures a ConversationsSelect.
type ConversationsSelectOption func(*ConversationsSelect)

// NewConversationsSelect creates a conversation select element.
func NewConversationsSelect(opts ...ConversationsSelectOption) ConversationsSelect {
	c := ConversationsSelect{}

	for _, opt := range opts {
		opt(&c)
	}

	return c
}

// WithConversationsSelectActionID sets the action_id.
func WithConversationsSelectActionID(id string) ConversationsSelectOption {
	return func(c *ConversationsSelect) {
		c.actionID = id
	}
}

// WithConversationsSelectPlaceholder sets placeholder text.
func WithConversationsSelectPlaceholder(text string) ConversationsSelectOption {
	return func(c *ConversationsSelect) {
		if pt, err := NewPlainText(text); err == nil {
			c.placeholder = &pt
		}
	}
}

// WithInitialConversation sets the initially selected conversation ID.
func WithInitialConversation(conversationID string) ConversationsSelectOption {
	return func(c *ConversationsSelect) {
		c.initialConversation = conversationID
	}
}

// WithDefaultToCurrentConversation defaults to current conversation.
func WithDefaultToCurrentConversation() ConversationsSelectOption {
	return func(c *ConversationsSelect) {
		c.defaultToCurrentConversation = true
	}
}

// WithConversationsFilter sets a filter for the conversation list.
func WithConversationsFilter(filter ConversationFilter) ConversationsSelectOption {
	return func(c *ConversationsSelect) {
		c.filter = &filter
	}
}

// WithConversationsSelectConfirm adds a confirmation dialog.
func WithConversationsSelectConfirm(confirm ConfirmDialog) ConversationsSelectOption {
	return func(c *ConversationsSelect) {
		c.confirm = &confirm
	}
}

// WithConversationsSelectFocusOnLoad sets auto-focus.
func WithConversationsSelectFocusOnLoad() ConversationsSelectOption {
	return func(c *ConversationsSelect) {
		c.focusOnLoad = true
	}
}

// WithResponseURLEnabled enables response URL for workflows.
func WithResponseURLEnabled() ConversationsSelectOption {
	return func(c *ConversationsSelect) {
		c.responseURLEnabled = true
	}
}

// MultiConversationsSelect allows selecting multiple conversations.
type MultiConversationsSelect struct {
	actionID                   string
	placeholder                *PlainText
	initialConversations       []string
	defaultToCurrentConversation bool
	filter                     *ConversationFilter
	maxSelectedItems           int
	confirm                    *ConfirmDialog
	focusOnLoad                bool
}

// Marker interface implementations
func (MultiConversationsSelect) sectionAccessory() {}
func (MultiConversationsSelect) actionsElement()   {}
func (MultiConversationsSelect) inputElement()     {}

// MarshalJSON implements json.Marshaler.
func (m MultiConversationsSelect) MarshalJSON() ([]byte, error) {
	out := map[string]any{
		"type": "multi_conversations_select",
	}
	if m.actionID != "" {
		out["action_id"] = m.actionID
	}
	if m.placeholder != nil {
		out["placeholder"] = m.placeholder
	}
	if len(m.initialConversations) > 0 {
		out["initial_conversations"] = m.initialConversations
	}
	if m.defaultToCurrentConversation {
		out["default_to_current_conversation"] = true
	}
	if m.filter != nil {
		out["filter"] = m.filter
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

// MultiConversationsSelectOption configures a MultiConversationsSelect.
type MultiConversationsSelectOption func(*MultiConversationsSelect)

// NewMultiConversationsSelect creates a multi-conversation select element.
func NewMultiConversationsSelect(opts ...MultiConversationsSelectOption) MultiConversationsSelect {
	m := MultiConversationsSelect{}

	for _, opt := range opts {
		opt(&m)
	}

	return m
}

// WithMultiConversationsSelectActionID sets the action_id.
func WithMultiConversationsSelectActionID(id string) MultiConversationsSelectOption {
	return func(m *MultiConversationsSelect) {
		m.actionID = id
	}
}

// WithMultiConversationsSelectPlaceholder sets placeholder text.
func WithMultiConversationsSelectPlaceholder(text string) MultiConversationsSelectOption {
	return func(m *MultiConversationsSelect) {
		if pt, err := NewPlainText(text); err == nil {
			m.placeholder = &pt
		}
	}
}

// WithInitialConversations sets initially selected conversation IDs.
func WithInitialConversations(conversationIDs ...string) MultiConversationsSelectOption {
	return func(m *MultiConversationsSelect) {
		m.initialConversations = conversationIDs
	}
}

// WithMultiConversationsDefaultToCurrentConversation defaults to current.
func WithMultiConversationsDefaultToCurrentConversation() MultiConversationsSelectOption {
	return func(m *MultiConversationsSelect) {
		m.defaultToCurrentConversation = true
	}
}

// WithMultiConversationsFilter sets a filter for the conversation list.
func WithMultiConversationsFilter(filter ConversationFilter) MultiConversationsSelectOption {
	return func(m *MultiConversationsSelect) {
		m.filter = &filter
	}
}

// WithMultiConversationsSelectMaxItems sets max selectable items.
func WithMultiConversationsSelectMaxItems(max int) MultiConversationsSelectOption {
	return func(m *MultiConversationsSelect) {
		m.maxSelectedItems = max
	}
}

// WithMultiConversationsSelectConfirm adds a confirmation dialog.
func WithMultiConversationsSelectConfirm(confirm ConfirmDialog) MultiConversationsSelectOption {
	return func(m *MultiConversationsSelect) {
		m.confirm = &confirm
	}
}

// WithMultiConversationsSelectFocusOnLoad sets auto-focus.
func WithMultiConversationsSelectFocusOnLoad() MultiConversationsSelectOption {
	return func(m *MultiConversationsSelect) {
		m.focusOnLoad = true
	}
}
