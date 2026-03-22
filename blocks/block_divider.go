package blocks

import "encoding/json"

// Divider is a visual separator between blocks.
type Divider struct {
	blockID string
}

// Marker interface implementation
func (Divider) block() {}

// MarshalJSON implements json.Marshaler.
func (d Divider) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type": "divider",
	}
	if d.blockID != "" {
		m["block_id"] = d.blockID
	}
	return json.Marshal(m)
}

// DividerOption configures a Divider.
type DividerOption func(*Divider)

// NewDivider creates a new divider block.
func NewDivider(opts ...DividerOption) Divider {
	d := Divider{}

	for _, opt := range opts {
		opt(&d)
	}

	return d
}

// WithDividerBlockID sets the block_id.
func WithDividerBlockID(id string) DividerOption {
	return func(d *Divider) {
		d.blockID = id
	}
}
