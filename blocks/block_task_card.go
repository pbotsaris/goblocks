package blocks

import "encoding/json"

// TaskCard displays a single task, representing a single action.
// Available in messages only.
type TaskCard struct {
	taskID      string
	title       PlainText
	description *PlainText
	status      TaskCardStatus
	sources     []URLSource
	blockID     string
}

// Marker interface implementation
func (TaskCard) block() {}

// MarshalJSON implements json.Marshaler.
func (t TaskCard) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type":    "task_card",
		"task_id": t.taskID,
		"title":   t.title,
		"status":  t.status,
	}
	if t.description != nil {
		m["description"] = t.description
	}
	if len(t.sources) > 0 {
		m["sources"] = t.sources
	}
	if t.blockID != "" {
		m["block_id"] = t.blockID
	}
	return json.Marshal(m)
}

// TaskCardOption configures a TaskCard block.
type TaskCardOption func(*TaskCard)

// NewTaskCard creates a new task card block.
func NewTaskCard(taskID, title string, status TaskCardStatus, opts ...TaskCardOption) (TaskCard, error) {
	if err := validateRequired("task_id", taskID); err != nil {
		return TaskCard{}, err
	}
	if err := validateRequired("title", title); err != nil {
		return TaskCard{}, err
	}

	pt, err := NewPlainText(title)
	if err != nil {
		return TaskCard{}, err
	}

	t := TaskCard{
		taskID: taskID,
		title:  pt,
		status: status,
	}

	for _, opt := range opts {
		opt(&t)
	}

	return t, nil
}

// MustTaskCard creates a TaskCard or panics on error.
func MustTaskCard(taskID, title string, status TaskCardStatus, opts ...TaskCardOption) TaskCard {
	t, err := NewTaskCard(taskID, title, status, opts...)
	if err != nil {
		panic(err)
	}
	return t
}

// WithTaskCardDescription sets the description.
func WithTaskCardDescription(text string) TaskCardOption {
	return func(t *TaskCard) {
		if pt, err := NewPlainText(text); err == nil {
			t.description = &pt
		}
	}
}

// WithTaskCardSources sets the URL sources.
func WithTaskCardSources(sources []URLSource) TaskCardOption {
	return func(t *TaskCard) {
		t.sources = sources
	}
}

// WithTaskCardBlockID sets the block_id.
func WithTaskCardBlockID(id string) TaskCardOption {
	return func(t *TaskCard) {
		t.blockID = id
	}
}

// TaskCardStatus represents the status of a task.
type TaskCardStatus string

const (
	TaskCardStatusOpen       TaskCardStatus = "open"
	TaskCardStatusInProgress TaskCardStatus = "in_progress"
	TaskCardStatusComplete   TaskCardStatus = "complete"
)

// URLSource represents a URL source for referencing within a task card.
type URLSource struct {
	url   string
	title *PlainText
}

// MarshalJSON implements json.Marshaler.
func (s URLSource) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type": "url_source",
		"url":  s.url,
	}
	if s.title != nil {
		m["title"] = s.title
	}
	return json.Marshal(m)
}

// URLSourceOption configures a URLSource.
type URLSourceOption func(*URLSource)

// NewURLSource creates a new URL source element.
func NewURLSource(url string, opts ...URLSourceOption) (URLSource, error) {
	if err := validateRequired("url", url); err != nil {
		return URLSource{}, err
	}

	s := URLSource{url: url}

	for _, opt := range opts {
		opt(&s)
	}

	return s, nil
}

// MustURLSource creates a URLSource or panics on error.
func MustURLSource(url string, opts ...URLSourceOption) URLSource {
	s, err := NewURLSource(url, opts...)
	if err != nil {
		panic(err)
	}
	return s
}

// WithURLSourceTitle sets the title for the URL source.
func WithURLSourceTitle(title string) URLSourceOption {
	return func(s *URLSource) {
		if pt, err := NewPlainText(title); err == nil {
			s.title = &pt
		}
	}
}
