package blocks

import "encoding/json"

// Plan displays a collection of related tasks.
// Available in messages only.
type Plan struct {
	title    PlainText
	sections []PlanSection
	blockID  string
}

// Marker interface implementation
func (Plan) block() {}

// MarshalJSON implements json.Marshaler.
func (p Plan) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type":     "plan",
		"title":    p.title,
		"sections": p.sections,
	}
	if p.blockID != "" {
		m["block_id"] = p.blockID
	}
	return json.Marshal(m)
}

// PlanOption configures a Plan block.
type PlanOption func(*Plan)

// NewPlan creates a new plan block.
func NewPlan(title string, sections []PlanSection, opts ...PlanOption) (Plan, error) {
	if err := validateRequired("title", title); err != nil {
		return Plan{}, err
	}
	if err := validateMinItems("sections", sections, 1); err != nil {
		return Plan{}, err
	}

	pt, err := NewPlainText(title)
	if err != nil {
		return Plan{}, err
	}

	p := Plan{
		title:    pt,
		sections: sections,
	}

	for _, opt := range opts {
		opt(&p)
	}

	return p, nil
}

// MustPlan creates a Plan or panics on error.
func MustPlan(title string, sections []PlanSection, opts ...PlanOption) Plan {
	p, err := NewPlan(title, sections, opts...)
	if err != nil {
		panic(err)
	}
	return p
}

// WithPlanBlockID sets the block_id.
func WithPlanBlockID(id string) PlanOption {
	return func(p *Plan) {
		p.blockID = id
	}
}

// PlanSection represents a section within a plan.
type PlanSection struct {
	title PlainText
	items []PlanItem
}

// MarshalJSON implements json.Marshaler.
func (s PlanSection) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"title": s.title,
		"items": s.items,
	})
}

// NewPlanSection creates a new plan section.
func NewPlanSection(title string, items []PlanItem) (PlanSection, error) {
	if err := validateRequired("title", title); err != nil {
		return PlanSection{}, err
	}
	if err := validateMinItems("items", items, 1); err != nil {
		return PlanSection{}, err
	}

	pt, err := NewPlainText(title)
	if err != nil {
		return PlanSection{}, err
	}

	return PlanSection{
		title: pt,
		items: items,
	}, nil
}

// MustPlanSection creates a PlanSection or panics on error.
func MustPlanSection(title string, items []PlanItem) PlanSection {
	s, err := NewPlanSection(title, items)
	if err != nil {
		panic(err)
	}
	return s
}

// PlanItem represents an item within a plan section.
type PlanItem struct {
	text   PlainText
	status PlanItemStatus
}

// MarshalJSON implements json.Marshaler.
func (i PlanItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"text":   i.text,
		"status": i.status,
	})
}

// PlanItemStatus represents the status of a plan item.
type PlanItemStatus string

const (
	PlanItemStatusPending    PlanItemStatus = "pending"
	PlanItemStatusInProgress PlanItemStatus = "in_progress"
	PlanItemStatusComplete   PlanItemStatus = "complete"
)

// NewPlanItem creates a new plan item.
func NewPlanItem(text string, status PlanItemStatus) (PlanItem, error) {
	if err := validateRequired("text", text); err != nil {
		return PlanItem{}, err
	}

	pt, err := NewPlainText(text)
	if err != nil {
		return PlanItem{}, err
	}

	return PlanItem{
		text:   pt,
		status: status,
	}, nil
}

// MustPlanItem creates a PlanItem or panics on error.
func MustPlanItem(text string, status PlanItemStatus) PlanItem {
	i, err := NewPlanItem(text, status)
	if err != nil {
		panic(err)
	}
	return i
}
