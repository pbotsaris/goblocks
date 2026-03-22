package blocks

import "encoding/json"

// ImageElement represents an image element for use in section/context blocks.
type ImageElement struct {
	imageURL  string
	altText   string
	slackFile *SlackFile
}

// Marker interface implementations
func (ImageElement) sectionAccessory() {}
func (ImageElement) contextElement()   {}

// MarshalJSON implements json.Marshaler.
func (i ImageElement) MarshalJSON() ([]byte, error) {
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
	return json.Marshal(m)
}

// ImageElementOption configures an ImageElement.
type ImageElementOption func(*ImageElement)

// NewImageElement creates an image element from a URL.
// imageURL max: 3000 chars, altText max: 2000 chars
func NewImageElement(imageURL, altText string) (ImageElement, error) {
	if err := validateRequiredMaxLen("image_url", imageURL, 3000); err != nil {
		return ImageElement{}, err
	}
	if err := validateRequiredMaxLen("alt_text", altText, 2000); err != nil {
		return ImageElement{}, err
	}

	return ImageElement{
		imageURL: imageURL,
		altText:  altText,
	}, nil
}

// NewImageElementFromSlackFile creates an image element from a Slack file.
func NewImageElementFromSlackFile(slackFile SlackFile, altText string) (ImageElement, error) {
	if err := validateRequiredMaxLen("alt_text", altText, 2000); err != nil {
		return ImageElement{}, err
	}

	return ImageElement{
		slackFile: &slackFile,
		altText:   altText,
	}, nil
}
