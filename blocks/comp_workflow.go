package blocks

import "encoding/json"

// InputParameter defines an input parameter for workflow triggers.
type InputParameter struct {
	name  string
	value string
}

// MarshalJSON implements json.Marshaler.
func (p InputParameter) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}{
		Name:  p.name,
		Value: p.value,
	})
}

// NewInputParameter creates a new input parameter.
func NewInputParameter(name, value string) (InputParameter, error) {
	if err := validateRequired("name", name); err != nil {
		return InputParameter{}, err
	}
	if err := validateRequired("value", value); err != nil {
		return InputParameter{}, err
	}
	return InputParameter{name: name, value: value}, nil
}

// MustInputParameter creates an InputParameter or panics on error.
func MustInputParameter(name, value string) InputParameter {
	p, err := NewInputParameter(name, value)
	if err != nil {
		panic(err)
	}
	return p
}

// Trigger defines a workflow trigger with URL and optional input parameters.
type Trigger struct {
	url                          string
	customizableInputParameters []InputParameter
}

// MarshalJSON implements json.Marshaler.
func (t Trigger) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"url": t.url,
	}
	if len(t.customizableInputParameters) > 0 {
		m["customizable_input_parameters"] = t.customizableInputParameters
	}
	return json.Marshal(m)
}

// TriggerOption configures a Trigger.
type TriggerOption func(*Trigger)

// NewTrigger creates a new trigger with the given URL.
func NewTrigger(url string, opts ...TriggerOption) (Trigger, error) {
	if err := validateRequired("url", url); err != nil {
		return Trigger{}, err
	}

	t := Trigger{url: url}
	for _, opt := range opts {
		opt(&t)
	}
	return t, nil
}

// MustTrigger creates a Trigger or panics on error.
func MustTrigger(url string, opts ...TriggerOption) Trigger {
	t, err := NewTrigger(url, opts...)
	if err != nil {
		panic(err)
	}
	return t
}

// WithInputParameters sets the customizable input parameters.
func WithInputParameters(params ...InputParameter) TriggerOption {
	return func(t *Trigger) {
		t.customizableInputParameters = params
	}
}

// Workflow defines a workflow object containing trigger information.
type Workflow struct {
	trigger Trigger
}

// MarshalJSON implements json.Marshaler.
func (w Workflow) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Trigger Trigger `json:"trigger"`
	}{
		Trigger: w.trigger,
	})
}

// NewWorkflow creates a new workflow with the given trigger.
func NewWorkflow(trigger Trigger) Workflow {
	return Workflow{trigger: trigger}
}
