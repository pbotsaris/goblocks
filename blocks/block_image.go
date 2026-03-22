package blocks

import "encoding/json"

// ImageBlock displays an image.
type ImageBlock struct {
	imageURL  string
	altText   string
	title     *PlainText
	blockID   string
	slackFile *SlackFile
}

// Marker interface implementation
func (ImageBlock) block() {}

// MarshalJSON implements json.Marshaler.
func (i ImageBlock) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type":     "image",
		"alt_text": i.altText,
	}
	if i.imageURL != "" {
		m["image_url"] = i.imageURL
	}
	if i.slackFile != nil {
		m["slack_file"] = i.slackFile
	}
	if i.title != nil {
		m["title"] = i.title
	}
	if i.blockID != "" {
		m["block_id"] = i.blockID
	}
	return json.Marshal(m)
}

// ImageBlockOption configures an ImageBlock.
type ImageBlockOption func(*ImageBlock)

// NewImageBlock creates an image block from a URL.
// imageURL max: 3000 chars, altText max: 2000 chars
func NewImageBlock(imageURL, altText string, opts ...ImageBlockOption) (ImageBlock, error) {
	if err := validateRequiredMaxLen("image_url", imageURL, 3000); err != nil {
		return ImageBlock{}, err
	}
	if err := validateRequiredMaxLen("alt_text", altText, 2000); err != nil {
		return ImageBlock{}, err
	}

	i := ImageBlock{
		imageURL: imageURL,
		altText:  altText,
	}

	for _, opt := range opts {
		opt(&i)
	}

	return i, nil
}

// NewImageBlockFromSlackFile creates an image block from a Slack file.
func NewImageBlockFromSlackFile(slackFile SlackFile, altText string, opts ...ImageBlockOption) (ImageBlock, error) {
	if err := validateRequiredMaxLen("alt_text", altText, 2000); err != nil {
		return ImageBlock{}, err
	}

	i := ImageBlock{
		slackFile: &slackFile,
		altText:   altText,
	}

	for _, opt := range opts {
		opt(&i)
	}

	return i, nil
}

// WithImageBlockTitle sets the image title.
// Max 2000 characters.
func WithImageBlockTitle(title string) ImageBlockOption {
	return func(i *ImageBlock) {
		if pt, err := NewPlainText(title); err == nil {
			i.title = &pt
		}
	}
}

// WithImageBlockID sets the block_id.
func WithImageBlockID(id string) ImageBlockOption {
	return func(i *ImageBlock) {
		i.blockID = id
	}
}
