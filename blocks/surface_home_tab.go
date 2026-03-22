package blocks

import "encoding/json"

// HomeTab represents the App Home tab surface.
type HomeTab struct {
	blocks          []Block
	privateMetadata string
	callbackID      string
	externalID      string
}

// MarshalJSON implements json.Marshaler.
func (h HomeTab) MarshalJSON() ([]byte, error) {
	out := map[string]any{
		"type":   "home",
		"blocks": h.blocks,
	}
	if h.privateMetadata != "" {
		out["private_metadata"] = h.privateMetadata
	}
	if h.callbackID != "" {
		out["callback_id"] = h.callbackID
	}
	if h.externalID != "" {
		out["external_id"] = h.externalID
	}
	return json.Marshal(out)
}

// HomeTabOption configures a HomeTab.
type HomeTabOption func(*HomeTab)

// NewHomeTab creates a new home tab.
// max 100 blocks
func NewHomeTab(blocks []Block, opts ...HomeTabOption) (HomeTab, error) {
	if err := validateMaxItems("blocks", blocks, 100); err != nil {
		return HomeTab{}, err
	}

	h := HomeTab{blocks: blocks}

	for _, opt := range opts {
		opt(&h)
	}

	return h, nil
}

// MustHomeTab creates a HomeTab or panics on error.
func MustHomeTab(blocks []Block, opts ...HomeTabOption) HomeTab {
	h, err := NewHomeTab(blocks, opts...)
	if err != nil {
		panic(err)
	}
	return h
}

// WithHomeTabPrivateMetadata sets private metadata.
// Max 3000 characters.
func WithHomeTabPrivateMetadata(metadata string) HomeTabOption {
	return func(h *HomeTab) {
		h.privateMetadata = metadata
	}
}

// WithHomeTabCallbackID sets the callback ID.
// Max 255 characters.
func WithHomeTabCallbackID(id string) HomeTabOption {
	return func(h *HomeTab) {
		h.callbackID = id
	}
}

// WithHomeTabExternalID sets an external ID.
func WithHomeTabExternalID(id string) HomeTabOption {
	return func(h *HomeTab) {
		h.externalID = id
	}
}
