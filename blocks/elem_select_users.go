package blocks

import "encoding/json"

// UsersSelect allows users to select a user from the workspace.
type UsersSelect struct {
	actionID    string
	placeholder *PlainText
	initialUser string
	confirm     *ConfirmDialog
	focusOnLoad bool
}

// Marker interface implementations
func (UsersSelect) sectionAccessory() {}
func (UsersSelect) actionsElement()   {}
func (UsersSelect) inputElement()     {}

// MarshalJSON implements json.Marshaler.
func (u UsersSelect) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type": "users_select",
	}
	if u.actionID != "" {
		m["action_id"] = u.actionID
	}
	if u.placeholder != nil {
		m["placeholder"] = u.placeholder
	}
	if u.initialUser != "" {
		m["initial_user"] = u.initialUser
	}
	if u.confirm != nil {
		m["confirm"] = u.confirm
	}
	if u.focusOnLoad {
		m["focus_on_load"] = true
	}
	return json.Marshal(m)
}

// UsersSelectOption configures a UsersSelect.
type UsersSelectOption func(*UsersSelect)

// NewUsersSelect creates a user select element.
func NewUsersSelect(opts ...UsersSelectOption) UsersSelect {
	u := UsersSelect{}

	for _, opt := range opts {
		opt(&u)
	}

	return u
}

// WithUsersSelectActionID sets the action_id.
func WithUsersSelectActionID(id string) UsersSelectOption {
	return func(u *UsersSelect) {
		u.actionID = id
	}
}

// WithUsersSelectPlaceholder sets placeholder text.
func WithUsersSelectPlaceholder(text string) UsersSelectOption {
	return func(u *UsersSelect) {
		if pt, err := NewPlainText(text); err == nil {
			u.placeholder = &pt
		}
	}
}

// WithInitialUser sets the initially selected user ID.
func WithInitialUser(userID string) UsersSelectOption {
	return func(u *UsersSelect) {
		u.initialUser = userID
	}
}

// WithUsersSelectConfirm adds a confirmation dialog.
func WithUsersSelectConfirm(confirm ConfirmDialog) UsersSelectOption {
	return func(u *UsersSelect) {
		u.confirm = &confirm
	}
}

// WithUsersSelectFocusOnLoad sets auto-focus.
func WithUsersSelectFocusOnLoad() UsersSelectOption {
	return func(u *UsersSelect) {
		u.focusOnLoad = true
	}
}

// MultiUsersSelect allows selecting multiple users.
type MultiUsersSelect struct {
	actionID         string
	placeholder      *PlainText
	initialUsers     []string
	maxSelectedItems int
	confirm          *ConfirmDialog
	focusOnLoad      bool
}

// Marker interface implementations
func (MultiUsersSelect) sectionAccessory() {}
func (MultiUsersSelect) actionsElement()   {}
func (MultiUsersSelect) inputElement()     {}

// MarshalJSON implements json.Marshaler.
func (m MultiUsersSelect) MarshalJSON() ([]byte, error) {
	out := map[string]any{
		"type": "multi_users_select",
	}
	if m.actionID != "" {
		out["action_id"] = m.actionID
	}
	if m.placeholder != nil {
		out["placeholder"] = m.placeholder
	}
	if len(m.initialUsers) > 0 {
		out["initial_users"] = m.initialUsers
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

// MultiUsersSelectOption configures a MultiUsersSelect.
type MultiUsersSelectOption func(*MultiUsersSelect)

// NewMultiUsersSelect creates a multi-user select element.
func NewMultiUsersSelect(opts ...MultiUsersSelectOption) MultiUsersSelect {
	m := MultiUsersSelect{}

	for _, opt := range opts {
		opt(&m)
	}

	return m
}

// WithMultiUsersSelectActionID sets the action_id.
func WithMultiUsersSelectActionID(id string) MultiUsersSelectOption {
	return func(m *MultiUsersSelect) {
		m.actionID = id
	}
}

// WithMultiUsersSelectPlaceholder sets placeholder text.
func WithMultiUsersSelectPlaceholder(text string) MultiUsersSelectOption {
	return func(m *MultiUsersSelect) {
		if pt, err := NewPlainText(text); err == nil {
			m.placeholder = &pt
		}
	}
}

// WithInitialUsers sets initially selected user IDs.
func WithInitialUsers(userIDs ...string) MultiUsersSelectOption {
	return func(m *MultiUsersSelect) {
		m.initialUsers = userIDs
	}
}

// WithMultiUsersSelectMaxItems sets max selectable items.
func WithMultiUsersSelectMaxItems(max int) MultiUsersSelectOption {
	return func(m *MultiUsersSelect) {
		m.maxSelectedItems = max
	}
}

// WithMultiUsersSelectConfirm adds a confirmation dialog.
func WithMultiUsersSelectConfirm(confirm ConfirmDialog) MultiUsersSelectOption {
	return func(m *MultiUsersSelect) {
		m.confirm = &confirm
	}
}

// WithMultiUsersSelectFocusOnLoad sets auto-focus.
func WithMultiUsersSelectFocusOnLoad() MultiUsersSelectOption {
	return func(m *MultiUsersSelect) {
		m.focusOnLoad = true
	}
}
