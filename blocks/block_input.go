package blocks

import "encoding/json"

// Input collects information from users via input elements.
type Input struct {
	label          PlainText
	element        InputElement
	blockID        string
	hint           *PlainText
	dispatchAction bool
	optional       bool
}

// Marker interface implementation
func (Input) block() {}

// MarshalJSON implements json.Marshaler.
func (i Input) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type":    "input",
		"label":   i.label,
		"element": i.element,
	}
	if i.blockID != "" {
		m["block_id"] = i.blockID
	}
	if i.hint != nil {
		m["hint"] = i.hint
	}
	if i.dispatchAction {
		m["dispatch_action"] = true
	}
	if i.optional {
		m["optional"] = true
	}
	return json.Marshal(m)
}

// InputOption configures an Input block.
type InputOption func(*Input)

// NewInput creates an input block with required label and element.
// label max: 2000 characters
func NewInput(label string, element InputElement, opts ...InputOption) (Input, error) {
	if err := validateRequiredMaxLen("label", label, 2000); err != nil {
		return Input{}, err
	}
	if element == nil {
		return Input{}, newValidationError("element", "cannot be nil", ErrMissingRequired)
	}

	pt, err := NewPlainText(label)
	if err != nil {
		return Input{}, err
	}

	inp := Input{
		label:   pt,
		element: element,
	}

	for _, opt := range opts {
		opt(&inp)
	}

	return inp, nil
}

// MustInput creates an Input or panics on error.
func MustInput(label string, element InputElement, opts ...InputOption) Input {
	i, err := NewInput(label, element, opts...)
	if err != nil {
		panic(err)
	}
	return i
}

// WithInputBlockID sets the block_id.
func WithInputBlockID(id string) InputOption {
	return func(i *Input) {
		i.blockID = id
	}
}

// WithInputHint adds a hint below the input.
// Max 2000 characters.
func WithInputHint(hint string) InputOption {
	return func(i *Input) {
		if pt, err := NewPlainText(hint); err == nil {
			i.hint = &pt
		}
	}
}

// WithInputDispatchAction enables block_actions dispatch.
func WithInputDispatchAction() InputOption {
	return func(i *Input) {
		i.dispatchAction = true
	}
}

// AsInputOptional marks the input as optional.
func AsInputOptional() InputOption {
	return func(i *Input) {
		i.optional = true
	}
}
