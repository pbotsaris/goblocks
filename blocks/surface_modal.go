package blocks

import "encoding/json"

// Modal represents a modal dialog surface.
type Modal struct {
	title           PlainText
	blocks          []Block
	submit          *PlainText
	close           *PlainText
	privateMetadata string
	callbackID      string
	clearOnClose    bool
	notifyOnClose   bool
	externalID      string
	submitDisabled  bool
}

// MarshalJSON implements json.Marshaler.
func (m Modal) MarshalJSON() ([]byte, error) {
	out := map[string]any{
		"type":   "modal",
		"title":  m.title,
		"blocks": m.blocks,
	}
	if m.submit != nil {
		out["submit"] = m.submit
	}
	if m.close != nil {
		out["close"] = m.close
	}
	if m.privateMetadata != "" {
		out["private_metadata"] = m.privateMetadata
	}
	if m.callbackID != "" {
		out["callback_id"] = m.callbackID
	}
	if m.clearOnClose {
		out["clear_on_close"] = true
	}
	if m.notifyOnClose {
		out["notify_on_close"] = true
	}
	if m.externalID != "" {
		out["external_id"] = m.externalID
	}
	if m.submitDisabled {
		out["submit_disabled"] = true
	}
	return json.Marshal(out)
}

// ModalOption configures a Modal.
type ModalOption func(*Modal)

// NewModal creates a new modal.
// title max: 24 characters, max 100 blocks
func NewModal(title string, blocks []Block, opts ...ModalOption) (Modal, error) {
	if err := validateRequiredMaxLen("title", title, 24); err != nil {
		return Modal{}, err
	}
	if err := validateMaxItems("blocks", blocks, 100); err != nil {
		return Modal{}, err
	}

	pt, err := NewPlainText(title)
	if err != nil {
		return Modal{}, err
	}

	m := Modal{
		title:  pt,
		blocks: blocks,
	}

	for _, opt := range opts {
		opt(&m)
	}

	return m, nil
}

// MustModal creates a Modal or panics on error.
func MustModal(title string, blocks []Block, opts ...ModalOption) Modal {
	m, err := NewModal(title, blocks, opts...)
	if err != nil {
		panic(err)
	}
	return m
}

// WithModalSubmit sets the submit button text.
// Max 24 characters.
func WithModalSubmit(text string) ModalOption {
	return func(m *Modal) {
		if pt, err := NewPlainText(text); err == nil {
			m.submit = &pt
		}
	}
}

// WithModalClose sets the close button text.
// Max 24 characters.
func WithModalClose(text string) ModalOption {
	return func(m *Modal) {
		if pt, err := NewPlainText(text); err == nil {
			m.close = &pt
		}
	}
}

// WithModalPrivateMetadata sets private metadata.
// Max 3000 characters.
func WithModalPrivateMetadata(metadata string) ModalOption {
	return func(m *Modal) {
		m.privateMetadata = metadata
	}
}

// WithModalCallbackID sets the callback ID.
// Max 255 characters.
func WithModalCallbackID(id string) ModalOption {
	return func(m *Modal) {
		m.callbackID = id
	}
}

// WithModalClearOnClose clears all views when closed.
func WithModalClearOnClose() ModalOption {
	return func(m *Modal) {
		m.clearOnClose = true
	}
}

// WithModalNotifyOnClose sends view_closed event when closed.
func WithModalNotifyOnClose() ModalOption {
	return func(m *Modal) {
		m.notifyOnClose = true
	}
}

// WithModalExternalID sets an external ID.
func WithModalExternalID(id string) ModalOption {
	return func(m *Modal) {
		m.externalID = id
	}
}

// WithModalSubmitDisabled disables the submit button.
func WithModalSubmitDisabled() ModalOption {
	return func(m *Modal) {
		m.submitDisabled = true
	}
}
