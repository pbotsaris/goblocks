package blocks

import "encoding/json"

// Option represents an option composition object.
// Used in select menus, checkboxes, radio buttons, and overflow menus.
type Option struct {
	text        TextObject
	value       string
	description TextObject
	url         string // only for overflow menu options
}

// MarshalJSON implements json.Marshaler.
func (o Option) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"text":  o.text,
		"value": o.value,
	}
	if o.description != nil {
		m["description"] = o.description
	}
	if o.url != "" {
		m["url"] = o.url
	}
	return json.Marshal(m)
}

// Value returns the option's value.
func (o Option) Value() string {
	return o.value
}

// OptionConfig configures an Option.
type OptionConfig func(*Option)

// NewOption creates a new Option with plain text.
// text max: 75 chars, value max: 150 chars
func NewOption(text, value string, opts ...OptionConfig) (Option, error) {
	if err := validateRequiredMaxLen("text", text, 75); err != nil {
		return Option{}, err
	}
	if err := validateRequiredMaxLen("value", value, 150); err != nil {
		return Option{}, err
	}

	pt, err := NewPlainText(text)
	if err != nil {
		return Option{}, err
	}

	o := Option{
		text:  pt,
		value: value,
	}

	for _, opt := range opts {
		opt(&o)
	}

	return o, nil
}

// NewOptionWithMarkdown creates an Option with markdown text.
// Only valid for checkboxes and radio buttons.
func NewOptionWithMarkdown(text, value string, opts ...OptionConfig) (Option, error) {
	if err := validateRequiredMaxLen("text", text, 75); err != nil {
		return Option{}, err
	}
	if err := validateRequiredMaxLen("value", value, 150); err != nil {
		return Option{}, err
	}

	md, err := NewMarkdown(text)
	if err != nil {
		return Option{}, err
	}

	o := Option{
		text:  md,
		value: value,
	}

	for _, opt := range opts {
		opt(&o)
	}

	return o, nil
}

// WithDescription adds a description to the option.
// Max 75 characters.
func WithDescription(description string) OptionConfig {
	return func(o *Option) {
		if description != "" {
			pt, _ := NewPlainText(description)
			o.description = pt
		}
	}
}

// WithMarkdownDescription adds a markdown description (for checkboxes/radio).
func WithMarkdownDescription(description string) OptionConfig {
	return func(o *Option) {
		if description != "" {
			md, _ := NewMarkdown(description)
			o.description = md
		}
	}
}

// WithOptionURL adds a URL to the option (only for overflow menu).
// Max 3000 characters.
func WithOptionURL(url string) OptionConfig {
	return func(o *Option) {
		o.url = url
	}
}

// OptionGroup represents an option group composition object.
// Groups options in select menus.
type OptionGroup struct {
	label   PlainText
	options []Option
}

// MarshalJSON implements json.Marshaler.
func (g OptionGroup) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Label   PlainText `json:"label"`
		Options []Option  `json:"options"`
	}{
		Label:   g.label,
		Options: g.options,
	})
}

// NewOptionGroup creates a new OptionGroup.
// label max: 75 chars, options max: 100 items
func NewOptionGroup(label string, options []Option) (OptionGroup, error) {
	if err := validateRequiredMaxLen("label", label, 75); err != nil {
		return OptionGroup{}, err
	}
	if err := validateMinItems("options", options, 1); err != nil {
		return OptionGroup{}, err
	}
	if err := validateMaxItems("options", options, 100); err != nil {
		return OptionGroup{}, err
	}

	pt, err := NewPlainText(label)
	if err != nil {
		return OptionGroup{}, err
	}

	return OptionGroup{
		label:   pt,
		options: options,
	}, nil
}
