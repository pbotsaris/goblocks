package blocks

import (
	"encoding/json"
	"testing"
)

func TestActions(t *testing.T) {
	t.Run("creates valid actions block with buttons", func(t *testing.T) {
		btn1, _ := NewButton("Button 1", WithActionID("btn1"))
		btn2, _ := NewButton("Button 2", WithActionID("btn2"))
		actions, err := NewActions([]ActionsElement{btn1, btn2})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(actions)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "actions" {
			t.Errorf("got type %v, want actions", result["type"])
		}

		elements := result["elements"].([]any)
		if len(elements) != 2 {
			t.Errorf("got %d elements, want 2", len(elements))
		}
	})

	t.Run("includes block_id when set", func(t *testing.T) {
		btn, _ := NewButton("Button")
		actions, _ := NewActions([]ActionsElement{btn}, WithActionsBlockID("actions_1"))

		data, _ := json.Marshal(actions)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["block_id"] != "actions_1" {
			t.Errorf("got block_id %v, want 'actions_1'", result["block_id"])
		}
	})

	t.Run("supports mixed element types", func(t *testing.T) {
		btn, _ := NewButton("Button")
		opt, _ := NewOption("Option", "opt")
		sel, _ := NewStaticSelect([]Option{opt})
		dp := NewDatePicker()
		actions, _ := NewActions([]ActionsElement{btn, sel, dp})

		data, _ := json.Marshal(actions)
		var result map[string]any
		mustUnmarshal(data, &result)

		elements := result["elements"].([]any)
		if len(elements) != 3 {
			t.Errorf("got %d elements, want 3", len(elements))
		}

		// Verify element types
		elem0 := elements[0].(map[string]any)
		elem1 := elements[1].(map[string]any)
		elem2 := elements[2].(map[string]any)

		if elem0["type"] != "button" {
			t.Errorf("element 0 type = %v, want button", elem0["type"])
		}
		if elem1["type"] != "static_select" {
			t.Errorf("element 1 type = %v, want static_select", elem1["type"])
		}
		if elem2["type"] != "datepicker" {
			t.Errorf("element 2 type = %v, want datepicker", elem2["type"])
		}
	})

	t.Run("rejects empty elements", func(t *testing.T) {
		_, err := NewActions([]ActionsElement{})
		if err == nil {
			t.Error("expected error for empty elements")
		}
	})

	t.Run("implements Block interface", func(t *testing.T) {
		btn, _ := NewButton("test")
		actions, _ := NewActions([]ActionsElement{btn})
		var _ Block = actions
	})
}
