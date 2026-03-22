package blocks

import (
	"encoding/json"
	"testing"
)

func TestOption(t *testing.T) {
	t.Run("creates valid option", func(t *testing.T) {
		opt, err := NewOption("Label", "value")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(opt)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		text := result["text"].(map[string]any)
		if text["text"] != "Label" {
			t.Errorf("got label %v, want Label", text["text"])
		}
		if result["value"] != "value" {
			t.Errorf("got value %v, want value", result["value"])
		}
	})

	t.Run("includes description when set", func(t *testing.T) {
		opt, err := NewOption("Label", "value", WithDescription("A description"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(opt)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		desc := result["description"].(map[string]any)
		if desc["text"] != "A description" {
			t.Errorf("got description %v, want 'A description'", desc["text"])
		}
	})

	t.Run("includes URL when set", func(t *testing.T) {
		opt, err := NewOption("Label", "value", WithOptionURL("https://example.com"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(opt)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["url"] != "https://example.com" {
			t.Errorf("got url %v, want https://example.com", result["url"])
		}
	})

	t.Run("rejects empty label", func(t *testing.T) {
		_, err := NewOption("", "value")
		if err == nil {
			t.Error("expected error for empty label")
		}
	})

	t.Run("rejects empty value", func(t *testing.T) {
		_, err := NewOption("Label", "")
		if err == nil {
			t.Error("expected error for empty value")
		}
	})
}

func TestOptionGroup(t *testing.T) {
	t.Run("creates valid option group", func(t *testing.T) {
		opt1, _ := NewOption("Option 1", "opt1")
		opt2, _ := NewOption("Option 2", "opt2")

		group, err := NewOptionGroup("Group", []Option{opt1, opt2})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(group)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		label := result["label"].(map[string]any)
		if label["text"] != "Group" {
			t.Errorf("got label %v, want Group", label["text"])
		}

		options := result["options"].([]any)
		if len(options) != 2 {
			t.Errorf("got %d options, want 2", len(options))
		}
	})

	t.Run("rejects empty label", func(t *testing.T) {
		opt, _ := NewOption("Option", "opt")
		_, err := NewOptionGroup("", []Option{opt})
		if err == nil {
			t.Error("expected error for empty label")
		}
	})

	t.Run("rejects empty options", func(t *testing.T) {
		_, err := NewOptionGroup("Group", []Option{})
		if err == nil {
			t.Error("expected error for empty options")
		}
	})
}
