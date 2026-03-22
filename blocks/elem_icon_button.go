package blocks

import "encoding/json"

// IconButton displays an icon button for performing actions.
// Used in context actions blocks.
type IconButton struct {
	icon     Icon
	actionID string
	altText  string
}

// Icon represents an icon for the icon button.
type Icon struct {
	name string
}

// MarshalJSON implements json.Marshaler for Icon.
func (i Icon) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"type"`
		Name string `json:"name,omitempty"`
	}{
		Type: "plain_icon",
		Name: i.name,
	})
}

// NewIcon creates a new icon with the given name.
func NewIcon(name string) Icon {
	return Icon{name: name}
}

// Marker interface implementation
func (IconButton) contextActionsElement() {}

// MarshalJSON implements json.Marshaler.
func (b IconButton) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type": "icon_button",
		"icon": b.icon,
	}
	if b.actionID != "" {
		m["action_id"] = b.actionID
	}
	if b.altText != "" {
		m["alt_text"] = b.altText
	}
	return json.Marshal(m)
}

// IconButtonOption configures an IconButton.
type IconButtonOption func(*IconButton)

// NewIconButton creates a new icon button element.
func NewIconButton(icon Icon, opts ...IconButtonOption) IconButton {
	b := IconButton{icon: icon}
	for _, opt := range opts {
		opt(&b)
	}
	return b
}

// WithIconButtonActionID sets the action_id.
func WithIconButtonActionID(id string) IconButtonOption {
	return func(b *IconButton) {
		b.actionID = id
	}
}

// WithIconButtonAltText sets the alt text for accessibility.
func WithIconButtonAltText(text string) IconButtonOption {
	return func(b *IconButton) {
		b.altText = text
	}
}
