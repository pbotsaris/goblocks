package blocks

import (
	"encoding/json"
	"testing"
)

func TestInput(t *testing.T) {
	t.Run("creates valid input block with plain text input", func(t *testing.T) {
		pti := NewPlainTextInput()
		input, err := NewInput("Enter text", pti)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(input)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "input" {
			t.Errorf("got type %v, want input", result["type"])
		}

		label := result["label"].(map[string]any)
		if label["text"] != "Enter text" {
			t.Errorf("got label %v, want 'Enter text'", label["text"])
		}

		element := result["element"].(map[string]any)
		if element["type"] != "plain_text_input" {
			t.Errorf("got element type %v, want 'plain_text_input'", element["type"])
		}
	})

	t.Run("includes block_id when set", func(t *testing.T) {
		pti := NewPlainTextInput()
		input, _ := NewInput("Label", pti, WithInputBlockID("input_1"))

		data, _ := json.Marshal(input)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["block_id"] != "input_1" {
			t.Errorf("got block_id %v, want 'input_1'", result["block_id"])
		}
	})

	t.Run("includes hint when set", func(t *testing.T) {
		pti := NewPlainTextInput()
		input, _ := NewInput("Label", pti, WithInputHint("Enter your name"))

		data, _ := json.Marshal(input)
		var result map[string]any
		mustUnmarshal(data, &result)

		hint := result["hint"].(map[string]any)
		if hint["text"] != "Enter your name" {
			t.Errorf("got hint %v, want 'Enter your name'", hint["text"])
		}
	})

	t.Run("includes dispatch_action when set", func(t *testing.T) {
		pti := NewPlainTextInput()
		input, _ := NewInput("Label", pti, WithInputDispatchAction())

		data, _ := json.Marshal(input)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["dispatch_action"] != true {
			t.Error("expected dispatch_action to be true")
		}
	})

	t.Run("supports static select element", func(t *testing.T) {
		opt, _ := NewOption("Option", "opt")
		sel, _ := NewStaticSelect([]Option{opt})
		input, _ := NewInput("Select option", sel)

		data, _ := json.Marshal(input)
		var result map[string]any
		mustUnmarshal(data, &result)

		element := result["element"].(map[string]any)
		if element["type"] != "static_select" {
			t.Errorf("got element type %v, want 'static_select'", element["type"])
		}
	})

	t.Run("supports date picker element", func(t *testing.T) {
		dp := NewDatePicker()
		input, _ := NewInput("Select date", dp)

		data, _ := json.Marshal(input)
		var result map[string]any
		mustUnmarshal(data, &result)

		element := result["element"].(map[string]any)
		if element["type"] != "datepicker" {
			t.Errorf("got element type %v, want 'datepicker'", element["type"])
		}
	})

	t.Run("rejects empty label", func(t *testing.T) {
		pti := NewPlainTextInput()
		_, err := NewInput("", pti)
		if err == nil {
			t.Error("expected error for empty label")
		}
	})

	t.Run("implements Block interface", func(t *testing.T) {
		pti := NewPlainTextInput()
		input, _ := NewInput("Label", pti)
		var _ Block = input
	})
}
