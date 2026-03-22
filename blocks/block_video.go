package blocks

import "encoding/json"

// Video displays an embedded video player.
type Video struct {
	altText         string
	title           PlainText
	thumbnailURL    string
	videoURL        string
	authorName      string
	blockID         string
	description     *PlainText
	providerIconURL string
	providerName    string
	titleURL        string
}

// Marker interface implementation
func (Video) block() {}

// MarshalJSON implements json.Marshaler.
func (v Video) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type":          "video",
		"alt_text":      v.altText,
		"title":         v.title,
		"thumbnail_url": v.thumbnailURL,
		"video_url":     v.videoURL,
	}
	if v.authorName != "" {
		m["author_name"] = v.authorName
	}
	if v.blockID != "" {
		m["block_id"] = v.blockID
	}
	if v.description != nil {
		m["description"] = v.description
	}
	if v.providerIconURL != "" {
		m["provider_icon_url"] = v.providerIconURL
	}
	if v.providerName != "" {
		m["provider_name"] = v.providerName
	}
	if v.titleURL != "" {
		m["title_url"] = v.titleURL
	}
	return json.Marshal(m)
}

// VideoOption configures a Video block.
type VideoOption func(*Video)

// NewVideo creates a new video block.
// altText: max 2000 characters
// title: max 200 characters
func NewVideo(altText, title, thumbnailURL, videoURL string, opts ...VideoOption) (Video, error) {
	if err := validateRequiredMaxLen("alt_text", altText, 2000); err != nil {
		return Video{}, err
	}
	if err := validateRequiredMaxLen("title", title, 200); err != nil {
		return Video{}, err
	}
	if err := validateRequired("thumbnail_url", thumbnailURL); err != nil {
		return Video{}, err
	}
	if err := validateRequired("video_url", videoURL); err != nil {
		return Video{}, err
	}

	pt, err := NewPlainText(title)
	if err != nil {
		return Video{}, err
	}

	v := Video{
		altText:      altText,
		title:        pt,
		thumbnailURL: thumbnailURL,
		videoURL:     videoURL,
	}

	for _, opt := range opts {
		opt(&v)
	}

	return v, nil
}

// MustVideo creates a Video or panics on error.
func MustVideo(altText, title, thumbnailURL, videoURL string, opts ...VideoOption) Video {
	v, err := NewVideo(altText, title, thumbnailURL, videoURL, opts...)
	if err != nil {
		panic(err)
	}
	return v
}

// WithVideoAuthorName sets the author name.
func WithVideoAuthorName(name string) VideoOption {
	return func(v *Video) {
		v.authorName = name
	}
}

// WithVideoBlockID sets the block_id.
func WithVideoBlockID(id string) VideoOption {
	return func(v *Video) {
		v.blockID = id
	}
}

// WithVideoDescription sets the description.
func WithVideoDescription(text string) VideoOption {
	return func(v *Video) {
		if pt, err := NewPlainText(text); err == nil {
			v.description = &pt
		}
	}
}

// WithVideoProviderIconURL sets the provider icon URL.
func WithVideoProviderIconURL(url string) VideoOption {
	return func(v *Video) {
		v.providerIconURL = url
	}
}

// WithVideoProviderName sets the provider name.
func WithVideoProviderName(name string) VideoOption {
	return func(v *Video) {
		v.providerName = name
	}
}

// WithVideoTitleURL sets the title URL (clickable link).
func WithVideoTitleURL(url string) VideoOption {
	return func(v *Video) {
		v.titleURL = url
	}
}
