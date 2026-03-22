package blocks

import "encoding/json"

// FileInput allows users to upload files.
type FileInput struct {
	actionID  string
	filetypes []string
	maxFiles  int
}

// Marker interface implementation
func (FileInput) inputElement() {}

// MarshalJSON implements json.Marshaler.
func (f FileInput) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type": "file_input",
	}
	if f.actionID != "" {
		m["action_id"] = f.actionID
	}
	if len(f.filetypes) > 0 {
		m["filetypes"] = f.filetypes
	}
	if f.maxFiles > 0 {
		m["max_files"] = f.maxFiles
	}
	return json.Marshal(m)
}

// FileInputOption configures a FileInput.
type FileInputOption func(*FileInput)

// NewFileInput creates a new file input element.
func NewFileInput(opts ...FileInputOption) FileInput {
	f := FileInput{}
	for _, opt := range opts {
		opt(&f)
	}
	return f
}

// WithFileInputActionID sets the action_id.
func WithFileInputActionID(id string) FileInputOption {
	return func(f *FileInput) {
		f.actionID = id
	}
}

// WithFileInputFiletypes sets the allowed file types.
// Example: []string{"pdf", "doc", "docx"}
func WithFileInputFiletypes(filetypes []string) FileInputOption {
	return func(f *FileInput) {
		f.filetypes = filetypes
	}
}

// WithFileInputMaxFiles sets the maximum number of files.
// Minimum 1, maximum 10. Defaults to 10 if not specified.
func WithFileInputMaxFiles(max int) FileInputOption {
	return func(f *FileInput) {
		if max >= 1 && max <= 10 {
			f.maxFiles = max
		}
	}
}
