package blocks

import (
	"encoding/json"
	"fmt"
)

// Builder provides a fluent API for building block collections.
type Builder struct {
	blocks []Block
	errors []error
}

// NewBuilder creates a new Builder.
func NewBuilder() *Builder {
	return &Builder{
		blocks: make([]Block, 0),
	}
}

// Add adds a block to the builder.
func (b *Builder) Add(block Block) *Builder {
	b.blocks = append(b.blocks, block)
	return b
}

// AddSection adds a section block with text.
func (b *Builder) AddSection(text TextObject, opts ...SectionOption) *Builder {
	section, err := NewSection(text, opts...)

	if err != nil {
		b.errors = append(b.errors, err)
		return b
	}

	b.blocks = append(b.blocks, section)

	return b
}

// AddSectionWithFields adds a section block with fields.
func (b *Builder) AddSectionWithFields(fields []TextObject, opts ...SectionOption) *Builder {
	section, err := NewSectionWithFields(fields, opts...)

	if err != nil {
		b.errors = append(b.errors, err)
		return b
	}

	b.blocks = append(b.blocks, section)

	return b
}

// AddDivider adds a divider block.
func (b *Builder) AddDivider(opts ...DividerOption) *Builder {
	b.blocks = append(b.blocks, NewDivider(opts...))

	return b
}

// AddHeader adds a header block.
func (b *Builder) AddHeader(text string, opts ...HeaderOption) *Builder {
	header, err := NewHeader(text, opts...)

	if err != nil {
		b.errors = append(b.errors, err)
		return b
	}

	b.blocks = append(b.blocks, header)

	return b
}

// AddActions adds an actions block.
func (b *Builder) AddActions(elements []ActionsElement, opts ...ActionsOption) *Builder {
	actions, err := NewActions(elements, opts...)

	if err != nil {
		b.errors = append(b.errors, err)
		return b
	}

	b.blocks = append(b.blocks, actions)
	return b
}

// AddContext adds a context block.
func (b *Builder) AddContext(elements []ContextElement, opts ...ContextOption) *Builder {
	ctx, err := NewContext(elements, opts...)

	if err != nil {
		b.errors = append(b.errors, err)
		return b
	}

	b.blocks = append(b.blocks, ctx)
	return b
}

// AddInput adds an input block.
func (b *Builder) AddInput(label string, element InputElement, opts ...InputOption) *Builder {
	input, err := NewInput(label, element, opts...)

	if err != nil {
		b.errors = append(b.errors, err)
		return b
	}

	b.blocks = append(b.blocks, input)

	return b
}

// AddImage adds an image block.
func (b *Builder) AddImage(imageURL, altText string, opts ...ImageBlockOption) *Builder {
	img, err := NewImageBlock(imageURL, altText, opts...)

	if err != nil {
		b.errors = append(b.errors, err)
		return b
	}

	b.blocks = append(b.blocks, img)
	return b
}

// Build returns the blocks and any accumulated errors.
func (b *Builder) Build() ([]Block, error) {

	if len(b.errors) > 0 {
		return nil, fmt.Errorf("builder errors: %v", b.errors)
	}

	return b.blocks, nil
}

// MustBuild returns blocks or panics on error.
func (b *Builder) MustBuild() []Block {

	blocks, err := b.Build()

	if err != nil {
		panic(err)
	}

	return blocks
}

// Errors returns any accumulated errors.
func (b *Builder) Errors() []error {
	return b.errors
}

// HasErrors returns true if there are any errors.
func (b *Builder) HasErrors() bool {
	return len(b.errors) > 0
}

// ToModal converts the blocks to a modal.
func (b *Builder) ToModal(title string, opts ...ModalOption) (Modal, error) {
	blocks, err := b.Build()

	if err != nil {
		return Modal{}, err
	}

	return NewModal(title, blocks, opts...)
}

// MustToModal converts to a modal or panics on error.
func (b *Builder) MustToModal(title string, opts ...ModalOption) Modal {
	modal, err := b.ToModal(title, opts...)

	if err != nil {
		panic(err)
	}

	return modal
}

// ToMessage converts the blocks to a message.
func (b *Builder) ToMessage(fallbackText string, opts ...MessageOption) (Message, error) {
	blocks, err := b.Build()

	if err != nil {
		return Message{}, err
	}

	return NewMessage(fallbackText, blocks, opts...)
}

// MustToMessage converts to a message or panics on error.
func (b *Builder) MustToMessage(fallbackText string, opts ...MessageOption) Message {
	msg, err := b.ToMessage(fallbackText, opts...)

	if err != nil {
		panic(err)
	}

	return msg
}

// ToHomeTab converts the blocks to a home tab.
func (b *Builder) ToHomeTab(opts ...HomeTabOption) (HomeTab, error) {
	blocks, err := b.Build()

	if err != nil {
		return HomeTab{}, err
	}

	return NewHomeTab(blocks, opts...)
}

// MustToHomeTab converts to a home tab or panics on error.
func (b *Builder) MustToHomeTab(opts ...HomeTabOption) HomeTab {
	home, err := b.ToHomeTab(opts...)

	if err != nil {
		panic(err)
	}

	return home
}

// JSON returns the JSON representation of the blocks.
func (b *Builder) JSON() ([]byte, error) {
	blocks, err := b.Build()

	if err != nil {
		return nil, err
	}

	return json.Marshal(map[string]any{"blocks": blocks})
}

// PrettyJSON returns indented JSON.
func (b *Builder) PrettyJSON() ([]byte, error) {
	blocks, err := b.Build()

	if err != nil {
		return nil, err
	}

	return json.MarshalIndent(map[string]any{"blocks": blocks}, "", "  ")
}

// BlocksJSON returns just the blocks array as JSON.
func (b *Builder) BlocksJSON() ([]byte, error) {
	blocks, err := b.Build()

	if err != nil {
		return nil, err
	}

	return json.Marshal(blocks)
}
