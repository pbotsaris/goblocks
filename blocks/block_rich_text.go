package blocks

import "encoding/json"

// RichText displays formatted, structured representation of text.
type RichText struct {
	elements []RichTextElement
	blockID  string
}

// Marker interface implementation
func (RichText) block() {}

// MarshalJSON implements json.Marshaler.
func (r RichText) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type":     "rich_text",
		"elements": r.elements,
	}
	if r.blockID != "" {
		m["block_id"] = r.blockID
	}
	return json.Marshal(m)
}

// RichTextOption configures a RichText block.
type RichTextOption func(*RichText)

// NewRichText creates a new rich text block.
func NewRichText(elements []RichTextElement, opts ...RichTextOption) (RichText, error) {
	if err := validateMinItems("elements", elements, 1); err != nil {
		return RichText{}, err
	}

	r := RichText{elements: elements}

	for _, opt := range opts {
		opt(&r)
	}

	return r, nil
}

// MustRichText creates a RichText or panics on error.
func MustRichText(elements []RichTextElement, opts ...RichTextOption) RichText {
	r, err := NewRichText(elements, opts...)
	if err != nil {
		panic(err)
	}
	return r
}

// WithRichTextBlockID sets the block_id.
func WithRichTextBlockID(id string) RichTextOption {
	return func(r *RichText) {
		r.blockID = id
	}
}

// RichTextSection represents a section of rich text content.
type RichTextSection struct {
	elements []RichTextSectionElement
}

// Marker interface implementation
func (RichTextSection) richTextElement() {}

// MarshalJSON implements json.Marshaler.
func (s RichTextSection) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":     "rich_text_section",
		"elements": s.elements,
	})
}

// NewRichTextSection creates a new rich text section.
func NewRichTextSection(elements []RichTextSectionElement) RichTextSection {
	return RichTextSection{elements: elements}
}

// RichTextList represents a list in rich text.
type RichTextList struct {
	style    string // "bullet" or "ordered"
	elements []RichTextSection
	indent   int
	offset   int
	border   int
}

// Marker interface implementation
func (RichTextList) richTextElement() {}

// MarshalJSON implements json.Marshaler.
func (l RichTextList) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type":     "rich_text_list",
		"style":    l.style,
		"elements": l.elements,
	}
	if l.indent > 0 {
		m["indent"] = l.indent
	}
	if l.offset > 0 {
		m["offset"] = l.offset
	}
	if l.border > 0 {
		m["border"] = l.border
	}
	return json.Marshal(m)
}

// RichTextListOption configures a RichTextList.
type RichTextListOption func(*RichTextList)

// NewRichTextList creates a new rich text list.
// style: "bullet" or "ordered"
func NewRichTextList(style string, elements []RichTextSection, opts ...RichTextListOption) RichTextList {
	l := RichTextList{
		style:    style,
		elements: elements,
	}
	for _, opt := range opts {
		opt(&l)
	}
	return l
}

// WithRichTextListIndent sets the indent level (0-8).
func WithRichTextListIndent(indent int) RichTextListOption {
	return func(l *RichTextList) {
		if indent >= 0 && indent <= 8 {
			l.indent = indent
		}
	}
}

// WithRichTextListOffset sets the offset for ordered lists.
func WithRichTextListOffset(offset int) RichTextListOption {
	return func(l *RichTextList) {
		l.offset = offset
	}
}

// WithRichTextListBorder sets the border (0 or 1).
func WithRichTextListBorder(border int) RichTextListOption {
	return func(l *RichTextList) {
		if border == 0 || border == 1 {
			l.border = border
		}
	}
}

// RichTextPreformatted represents preformatted text (code block).
type RichTextPreformatted struct {
	elements []RichTextSectionElement
	border   int
}

// Marker interface implementation
func (RichTextPreformatted) richTextElement() {}

// MarshalJSON implements json.Marshaler.
func (p RichTextPreformatted) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type":     "rich_text_preformatted",
		"elements": p.elements,
	}
	if p.border > 0 {
		m["border"] = p.border
	}
	return json.Marshal(m)
}

// NewRichTextPreformatted creates a new preformatted text element.
func NewRichTextPreformatted(elements []RichTextSectionElement, border int) RichTextPreformatted {
	return RichTextPreformatted{
		elements: elements,
		border:   border,
	}
}

// RichTextQuote represents a quote in rich text.
type RichTextQuote struct {
	elements []RichTextSectionElement
	border   int
}

// Marker interface implementation
func (RichTextQuote) richTextElement() {}

// MarshalJSON implements json.Marshaler.
func (q RichTextQuote) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type":     "rich_text_quote",
		"elements": q.elements,
	}
	if q.border > 0 {
		m["border"] = q.border
	}
	return json.Marshal(m)
}

// NewRichTextQuote creates a new quote element.
func NewRichTextQuote(elements []RichTextSectionElement, border int) RichTextQuote {
	return RichTextQuote{
		elements: elements,
		border:   border,
	}
}

// RichTextSectionElement is an element within a rich text section.
type RichTextSectionElement interface {
	json.Marshaler
	richTextSectionElement()
}

// RichTextText represents plain or styled text in a rich text section.
type RichTextText struct {
	text   string
	style  *RichTextStyle
}

func (RichTextText) richTextSectionElement() {}

// MarshalJSON implements json.Marshaler.
func (t RichTextText) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type": "text",
		"text": t.text,
	}
	if t.style != nil {
		m["style"] = t.style
	}
	return json.Marshal(m)
}

// NewRichTextText creates a new text element for rich text sections.
func NewRichTextText(text string, style *RichTextStyle) RichTextText {
	return RichTextText{text: text, style: style}
}

// RichTextLink represents a link in a rich text section.
type RichTextLink struct {
	url   string
	text  string
	style *RichTextStyle
}

func (RichTextLink) richTextSectionElement() {}

// MarshalJSON implements json.Marshaler.
func (l RichTextLink) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type": "link",
		"url":  l.url,
	}
	if l.text != "" {
		m["text"] = l.text
	}
	if l.style != nil {
		m["style"] = l.style
	}
	return json.Marshal(m)
}

// NewRichTextLink creates a new link element for rich text sections.
func NewRichTextLink(url string, text string, style *RichTextStyle) RichTextLink {
	return RichTextLink{url: url, text: text, style: style}
}

// RichTextEmoji represents an emoji in a rich text section.
type RichTextEmoji struct {
	name string
}

func (RichTextEmoji) richTextSectionElement() {}

// MarshalJSON implements json.Marshaler.
func (e RichTextEmoji) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type": "emoji",
		"name": e.name,
	})
}

// NewRichTextEmoji creates a new emoji element for rich text sections.
func NewRichTextEmoji(name string) RichTextEmoji {
	return RichTextEmoji{name: name}
}

// RichTextUser represents a user mention in a rich text section.
type RichTextUser struct {
	userID string
	style  *RichTextStyle
}

func (RichTextUser) richTextSectionElement() {}

// MarshalJSON implements json.Marshaler.
func (u RichTextUser) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type":    "user",
		"user_id": u.userID,
	}
	if u.style != nil {
		m["style"] = u.style
	}
	return json.Marshal(m)
}

// NewRichTextUser creates a new user mention for rich text sections.
func NewRichTextUser(userID string, style *RichTextStyle) RichTextUser {
	return RichTextUser{userID: userID, style: style}
}

// RichTextChannel represents a channel mention in a rich text section.
type RichTextChannel struct {
	channelID string
	style     *RichTextStyle
}

func (RichTextChannel) richTextSectionElement() {}

// MarshalJSON implements json.Marshaler.
func (c RichTextChannel) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type":       "channel",
		"channel_id": c.channelID,
	}
	if c.style != nil {
		m["style"] = c.style
	}
	return json.Marshal(m)
}

// NewRichTextChannel creates a new channel mention for rich text sections.
func NewRichTextChannel(channelID string, style *RichTextStyle) RichTextChannel {
	return RichTextChannel{channelID: channelID, style: style}
}

// RichTextUserGroup represents a user group mention in a rich text section.
type RichTextUserGroup struct {
	usergroupID string
	style       *RichTextStyle
}

func (RichTextUserGroup) richTextSectionElement() {}

// MarshalJSON implements json.Marshaler.
func (g RichTextUserGroup) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type":         "usergroup",
		"usergroup_id": g.usergroupID,
	}
	if g.style != nil {
		m["style"] = g.style
	}
	return json.Marshal(m)
}

// NewRichTextUserGroup creates a new user group mention for rich text sections.
func NewRichTextUserGroup(usergroupID string, style *RichTextStyle) RichTextUserGroup {
	return RichTextUserGroup{usergroupID: usergroupID, style: style}
}

// RichTextBroadcast represents a broadcast mention (@channel, @here, @everyone).
type RichTextBroadcast struct {
	rang string // "channel", "here", or "everyone"
}

func (RichTextBroadcast) richTextSectionElement() {}

// MarshalJSON implements json.Marshaler.
func (b RichTextBroadcast) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":  "broadcast",
		"range": b.rang,
	})
}

// NewRichTextBroadcast creates a new broadcast mention.
// rang: "channel", "here", or "everyone"
func NewRichTextBroadcast(rang string) RichTextBroadcast {
	return RichTextBroadcast{rang: rang}
}

// RichTextStyle defines text styling options.
type RichTextStyle struct {
	Bold   bool `json:"bold,omitempty"`
	Italic bool `json:"italic,omitempty"`
	Strike bool `json:"strike,omitempty"`
	Code   bool `json:"code,omitempty"`
}

// NewRichTextStyle creates a new style configuration.
func NewRichTextStyle(bold, italic, strike, code bool) *RichTextStyle {
	return &RichTextStyle{
		Bold:   bold,
		Italic: italic,
		Strike: strike,
		Code:   code,
	}
}
