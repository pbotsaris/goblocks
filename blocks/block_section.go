package blocks

import "encoding/json"

// Section is a flexible block for displaying text and accessories.
type Section struct {
	text      TextObject
	blockID   string
	fields    []TextObject
	accessory SectionAccessory
}

// Marker interface implementation
func (Section) block() {}

// MarshalJSON implements json.Marshaler.
func (s Section) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type": "section",
	}
	if s.text != nil {
		m["text"] = s.text
	}
	if s.blockID != "" {
		m["block_id"] = s.blockID
	}
	if len(s.fields) > 0 {
		m["fields"] = s.fields
	}
	if s.accessory != nil {
		m["accessory"] = s.accessory
	}
	return json.Marshal(m)
}

// SectionOption configures a Section.
type SectionOption func(*Section)

// NewSection creates a section with text.
func NewSection(text TextObject, opts ...SectionOption) (Section, error) {
	if text == nil {
		return Section{}, newValidationError("text", "cannot be nil", ErrMissingRequired)
	}

	s := Section{text: text}

	for _, opt := range opts {
		opt(&s)
	}

	// Validate fields count if set
	if len(s.fields) > 10 {
		return Section{}, newValidationError("fields", "exceeds maximum of 10", ErrExceedsMaxItems)
	}

	return s, nil
}

// NewSectionWithFields creates a section with fields instead of text.
// Max 10 fields.
func NewSectionWithFields(fields []TextObject, opts ...SectionOption) (Section, error) {
	if err := validateMinItems("fields", fields, 1); err != nil {
		return Section{}, err
	}
	if err := validateMaxItems("fields", fields, 10); err != nil {
		return Section{}, err
	}

	s := Section{fields: fields}

	for _, opt := range opts {
		opt(&s)
	}

	return s, nil
}

// MustSection creates a Section or panics on error.
func MustSection(text TextObject, opts ...SectionOption) Section {
	s, err := NewSection(text, opts...)
	if err != nil {
		panic(err)
	}
	return s
}

// WithSectionBlockID sets the block_id.
func WithSectionBlockID(id string) SectionOption {
	return func(s *Section) {
		s.blockID = id
	}
}

// WithSectionFields adds fields to the section.
// Max 10 fields.
func WithSectionFields(fields ...TextObject) SectionOption {
	return func(s *Section) {
		s.fields = fields
	}
}

// WithSectionAccessory adds an accessory element.
// Type-safety is enforced via the SectionAccessory interface.
func WithSectionAccessory(accessory SectionAccessory) SectionOption {
	return func(s *Section) {
		s.accessory = accessory
	}
}
