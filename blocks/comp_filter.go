package blocks

import "encoding/json"

// ConversationType represents types of conversations to include in filter.
type ConversationType string

const (
	ConversationIM      ConversationType = "im"
	ConversationMPIM    ConversationType = "mpim"
	ConversationPrivate ConversationType = "private"
	ConversationPublic  ConversationType = "public"
)

// ConversationFilter provides filtering for conversation selector menus.
type ConversationFilter struct {
	include                       []ConversationType
	excludeExternalSharedChannels bool
	excludeBotUsers               bool
}

// MarshalJSON implements json.Marshaler.
func (f ConversationFilter) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)

	if len(f.include) > 0 {
		m["include"] = f.include
	}
	if f.excludeExternalSharedChannels {
		m["exclude_external_shared_channels"] = true
	}
	if f.excludeBotUsers {
		m["exclude_bot_users"] = true
	}
	return json.Marshal(m)
}

// ConversationFilterOption configures a ConversationFilter.
type ConversationFilterOption func(*ConversationFilter)

// NewConversationFilter creates a new conversation filter.
// At least one option must be set.
func NewConversationFilter(opts ...ConversationFilterOption) (ConversationFilter, error) {
	f := ConversationFilter{}

	for _, opt := range opts {
		opt(&f)
	}

	// Validate at least one field is set
	if len(f.include) == 0 && !f.excludeExternalSharedChannels && !f.excludeBotUsers {
		return ConversationFilter{}, newValidationError("filter", "at least one field must be set", ErrMissingRequired)
	}

	return f, nil
}

// WithInclude specifies which conversation types to include.
func WithInclude(types ...ConversationType) ConversationFilterOption {
	return func(f *ConversationFilter) {
		f.include = types
	}
}

// WithExcludeExternalSharedChannels excludes external shared channels.
func WithExcludeExternalSharedChannels() ConversationFilterOption {
	return func(f *ConversationFilter) {
		f.excludeExternalSharedChannels = true
	}
}

// WithExcludeBotUsers excludes bot users.
func WithExcludeBotUsers() ConversationFilterOption {
	return func(f *ConversationFilter) {
		f.excludeBotUsers = true
	}
}

// TriggerAction represents when to dispatch a block_actions payload.
type TriggerAction string

const (
	TriggerOnEnterPressed     TriggerAction = "on_enter_pressed"
	TriggerOnCharacterEntered TriggerAction = "on_character_entered"
)

// DispatchActionConfig defines when a plain-text input dispatches actions.
type DispatchActionConfig struct {
	triggerActionsOn []TriggerAction
}

// MarshalJSON implements json.Marshaler.
func (d DispatchActionConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		TriggerActionsOn []TriggerAction `json:"trigger_actions_on,omitempty"`
	}{
		TriggerActionsOn: d.triggerActionsOn,
	})
}

// NewDispatchActionConfig creates a new dispatch action configuration.
func NewDispatchActionConfig(triggers ...TriggerAction) DispatchActionConfig {
	return DispatchActionConfig{
		triggerActionsOn: triggers,
	}
}
