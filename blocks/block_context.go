package blocks

import "encoding/json"

// Context displays contextual information.
// Elements can be images and text (PlainText or Markdown).
type Context struct {
	elements []ContextElement
	blockID  string
}

// Marker interface implementation
func (Context) block() {}

// MarshalJSON implements json.Marshaler.
func (c Context) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type":     "context",
		"elements": c.elements,
	}
	if c.blockID != "" {
		m["block_id"] = c.blockID
	}
	return json.Marshal(m)
}

// ContextOption configures a Context.
type ContextOption func(*Context)

// NewContext creates a new context block.
// Max 10 elements.
func NewContext(elements []ContextElement, opts ...ContextOption) (Context, error) {
	if err := validateMinItems("elements", elements, 1); err != nil {
		return Context{}, err
	}
	if err := validateMaxItems("elements", elements, 10); err != nil {
		return Context{}, err
	}

	c := Context{elements: elements}

	for _, opt := range opts {
		opt(&c)
	}

	return c, nil
}

// WithContextBlockID sets the block_id.
func WithContextBlockID(id string) ContextOption {
	return func(c *Context) {
		c.blockID = id
	}
}
