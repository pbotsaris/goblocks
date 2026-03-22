package blocks

import "encoding/json"

// SlackFile represents a Slack file composition object.
// Used in image blocks/elements to reference Slack-hosted files.
// Either URL or ID must be provided, but not both.
type SlackFile struct {
	url string
	id  string
}

// MarshalJSON implements json.Marshaler.
func (s SlackFile) MarshalJSON() ([]byte, error) {
	m := make(map[string]string)
	if s.url != "" {
		m["url"] = s.url
	}
	if s.id != "" {
		m["id"] = s.id
	}
	return json.Marshal(m)
}

// NewSlackFileFromURL creates a SlackFile from a URL.
// The URL can be url_private or permalink of a Slack file.
func NewSlackFileFromURL(url string) (SlackFile, error) {
	if err := validateRequired("url", url); err != nil {
		return SlackFile{}, err
	}
	if err := validateMaxLen("url", url, 3000); err != nil {
		return SlackFile{}, err
	}
	return SlackFile{url: url}, nil
}

// NewSlackFileFromID creates a SlackFile from a Slack file ID.
func NewSlackFileFromID(id string) (SlackFile, error) {
	if err := validateRequired("id", id); err != nil {
		return SlackFile{}, err
	}
	return SlackFile{id: id}, nil
}
