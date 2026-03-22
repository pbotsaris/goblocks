package blocks

import "encoding/json"

// WorkflowButton allows users to run a link trigger with customizable inputs.
type WorkflowButton struct {
	text               PlainText
	workflow           Workflow
	actionID           string
	style              ButtonStyle
	accessibilityLabel string
}

// Marker interface implementations
func (WorkflowButton) sectionAccessory() {}
func (WorkflowButton) actionsElement()   {}

// MarshalJSON implements json.Marshaler.
func (w WorkflowButton) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type":     "workflow_button",
		"text":     w.text,
		"workflow": w.workflow,
	}
	if w.actionID != "" {
		m["action_id"] = w.actionID
	}
	if w.style != "" {
		m["style"] = w.style
	}
	if w.accessibilityLabel != "" {
		m["accessibility_label"] = w.accessibilityLabel
	}
	return json.Marshal(m)
}

// WorkflowButtonOption configures a WorkflowButton.
type WorkflowButtonOption func(*WorkflowButton)

// NewWorkflowButton creates a new workflow button element.
// text max: 75 characters
func NewWorkflowButton(text string, workflow Workflow, opts ...WorkflowButtonOption) (WorkflowButton, error) {
	if err := validateRequiredMaxLen("text", text, 75); err != nil {
		return WorkflowButton{}, err
	}

	pt, err := NewPlainText(text)
	if err != nil {
		return WorkflowButton{}, err
	}

	w := WorkflowButton{
		text:     pt,
		workflow: workflow,
	}

	for _, opt := range opts {
		opt(&w)
	}

	return w, nil
}

// MustWorkflowButton creates a WorkflowButton or panics on error.
func MustWorkflowButton(text string, workflow Workflow, opts ...WorkflowButtonOption) WorkflowButton {
	w, err := NewWorkflowButton(text, workflow, opts...)
	if err != nil {
		panic(err)
	}
	return w
}

// WithWorkflowButtonActionID sets the action_id.
func WithWorkflowButtonActionID(id string) WorkflowButtonOption {
	return func(w *WorkflowButton) {
		w.actionID = id
	}
}

// WithWorkflowButtonStyle sets the button style (primary or danger).
func WithWorkflowButtonStyle(style ButtonStyle) WorkflowButtonOption {
	return func(w *WorkflowButton) {
		w.style = style
	}
}

// WithWorkflowButtonAccessibilityLabel sets the accessibility label.
// Max 75 characters.
func WithWorkflowButtonAccessibilityLabel(label string) WorkflowButtonOption {
	return func(w *WorkflowButton) {
		if len(label) <= 75 {
			w.accessibilityLabel = label
		}
	}
}
