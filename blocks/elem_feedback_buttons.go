package blocks

import "encoding/json"

// FeedbackButtons displays buttons to indicate positive or negative feedback.
// Used in context actions blocks for AI/assistant responses.
type FeedbackButtons struct {
	actionID string
}

// Marker interface implementation
func (FeedbackButtons) contextActionsElement() {}

// MarshalJSON implements json.Marshaler.
func (f FeedbackButtons) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type": "feedback_buttons",
	}
	if f.actionID != "" {
		m["action_id"] = f.actionID
	}
	return json.Marshal(m)
}

// FeedbackButtonsOption configures FeedbackButtons.
type FeedbackButtonsOption func(*FeedbackButtons)

// NewFeedbackButtons creates a new feedback buttons element.
func NewFeedbackButtons(opts ...FeedbackButtonsOption) FeedbackButtons {
	f := FeedbackButtons{}
	for _, opt := range opts {
		opt(&f)
	}
	return f
}

// WithFeedbackButtonsActionID sets the action_id.
func WithFeedbackButtonsActionID(id string) FeedbackButtonsOption {
	return func(f *FeedbackButtons) {
		f.actionID = id
	}
}
