package blocks

import "encoding/json"

// Block is the base interface for all top-level layout blocks.
// The unexported method seals the interface - only types in this
// package can implement it.
type Block interface {
	json.Marshaler
	block()
}

// TextObject represents any text composition object (plain_text or mrkdwn).
type TextObject interface {
	json.Marshaler
	textObject()
}

// PlainTextOnly restricts to plain_text only.
// Used for headers, placeholders, button text, etc.
type PlainTextOnly interface {
	TextObject
	plainTextOnly()
}

// SectionAccessory marks elements valid as section block accessories.
// Valid types: Button, ImageElement, StaticSelect, MultiStaticSelect,
// ExternalSelect, MultiExternalSelect, UsersSelect, MultiUsersSelect,
// ConversationsSelect, MultiConversationsSelect, ChannelsSelect,
// MultiChannelsSelect, DatePicker, TimePicker, Checkboxes, RadioButtons, Overflow
type SectionAccessory interface {
	json.Marshaler
	sectionAccessory()
}

// ActionsElement marks elements valid in actions blocks.
// Valid types: Button, StaticSelect, MultiStaticSelect, ExternalSelect,
// MultiExternalSelect, UsersSelect, MultiUsersSelect, ConversationsSelect,
// MultiConversationsSelect, ChannelsSelect, MultiChannelsSelect,
// DatePicker, TimePicker, Checkboxes, RadioButtons, Overflow
type ActionsElement interface {
	json.Marshaler
	actionsElement()
}

// ContextElement marks elements valid in context blocks.
// Valid types: ImageElement, PlainText, Markdown
type ContextElement interface {
	json.Marshaler
	contextElement()
}

// InputElement marks elements valid in input blocks.
// Valid types: PlainTextInput, StaticSelect, MultiStaticSelect,
// ExternalSelect, MultiExternalSelect, UsersSelect, MultiUsersSelect,
// ConversationsSelect, MultiConversationsSelect, ChannelsSelect,
// MultiChannelsSelect, DatePicker, TimePicker, DatetimePicker,
// Checkboxes, RadioButtons, EmailInput, NumberInput, URLInput,
// FileInput, RichTextInput
type InputElement interface {
	json.Marshaler
	inputElement()
}

// RichTextElement marks elements valid in rich text blocks.
// Valid types: RichTextSection, RichTextList, RichTextPreformatted, RichTextQuote
type RichTextElement interface {
	json.Marshaler
	richTextElement()
}

// ContextActionsElement marks elements valid in context actions blocks.
// Valid types: FeedbackButtons, IconButton
type ContextActionsElement interface {
	json.Marshaler
	contextActionsElement()
}
