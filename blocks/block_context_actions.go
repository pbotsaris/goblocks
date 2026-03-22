package blocks

import "encoding/json"

// ContextActions displays actions as contextual information.
// Available in messages only. Contains feedback buttons and icon buttons.
type ContextActions struct {
	elements []ContextActionsElement
	blockID  string
}

// Marker interface implementation
func (ContextActions) block() {}

// MarshalJSON implements json.Marshaler.
func (c ContextActions) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type":     "context_actions",
		"elements": c.elements,
	}
	if c.blockID != "" {
		m["block_id"] = c.blockID
	}
	return json.Marshal(m)
}

// ContextActionsOption configures a ContextActions block.
type ContextActionsOption func(*ContextActions)

// NewContextActions creates a new context actions block.
// elements: 1-5 elements (FeedbackButtons, IconButton)
func NewContextActions(elements []ContextActionsElement, opts ...ContextActionsOption) (ContextActions, error) {
	if err := validateMinItems("elements", elements, 1); err != nil {
		return ContextActions{}, err
	}
	if err := validateMaxItems("elements", elements, 5); err != nil {
		return ContextActions{}, err
	}

	c := ContextActions{elements: elements}

	for _, opt := range opts {
		opt(&c)
	}

	return c, nil
}

// MustContextActions creates a ContextActions or panics on error.
func MustContextActions(elements []ContextActionsElement, opts ...ContextActionsOption) ContextActions {
	c, err := NewContextActions(elements, opts...)
	if err != nil {
		panic(err)
	}
	return c
}

// WithContextActionsBlockID sets the block_id.
func WithContextActionsBlockID(id string) ContextActionsOption {
	return func(c *ContextActions) {
		c.blockID = id
	}
}
