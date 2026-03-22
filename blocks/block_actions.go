package blocks

import "encoding/json"

// Actions holds multiple interactive elements.
type Actions struct {
	elements []ActionsElement
	blockID  string
}

// Marker interface implementation
func (Actions) block() {}

// MarshalJSON implements json.Marshaler.
func (a Actions) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type":     "actions",
		"elements": a.elements,
	}
	if a.blockID != "" {
		m["block_id"] = a.blockID
	}
	return json.Marshal(m)
}

// ActionsOption configures an Actions block.
type ActionsOption func(*Actions)

// NewActions creates an actions block.
// Max 25 elements.
func NewActions(elements []ActionsElement, opts ...ActionsOption) (Actions, error) {
	if err := validateMinItems("elements", elements, 1); err != nil {
		return Actions{}, err
	}
	if err := validateMaxItems("elements", elements, 25); err != nil {
		return Actions{}, err
	}

	a := Actions{elements: elements}

	for _, opt := range opts {
		opt(&a)
	}

	return a, nil
}

// WithActionsBlockID sets the block_id.
func WithActionsBlockID(id string) ActionsOption {
	return func(a *Actions) {
		a.blockID = id
	}
}
