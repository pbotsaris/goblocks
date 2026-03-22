package blocks

import "encoding/json"

// File displays information about a remote file.
// This block appears when retrieving messages that contain remote files.
// You cannot add this block directly to messages.
type File struct {
	externalID string
	source     string
	blockID    string
}

// Marker interface implementation
func (File) block() {}

// MarshalJSON implements json.Marshaler.
func (f File) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type":        "file",
		"external_id": f.externalID,
		"source":      f.source,
	}
	if f.blockID != "" {
		m["block_id"] = f.blockID
	}
	return json.Marshal(m)
}

// FileOption configures a File block.
type FileOption func(*File)

// NewFile creates a new file block.
// externalID is the external unique ID for the file.
// source is always "remote" for remote files.
func NewFile(externalID string, opts ...FileOption) (File, error) {
	if err := validateRequired("external_id", externalID); err != nil {
		return File{}, err
	}

	f := File{
		externalID: externalID,
		source:     "remote",
	}

	for _, opt := range opts {
		opt(&f)
	}

	return f, nil
}

// MustFile creates a File or panics on error.
func MustFile(externalID string, opts ...FileOption) File {
	f, err := NewFile(externalID, opts...)
	if err != nil {
		panic(err)
	}
	return f
}

// WithFileBlockID sets the block_id.
func WithFileBlockID(id string) FileOption {
	return func(f *File) {
		f.blockID = id
	}
}
