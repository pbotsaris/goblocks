package blocks

import (
	"encoding/json"
	"testing"
)

func TestSection(t *testing.T) {
	t.Run("creates valid section with text", func(t *testing.T) {
		text, _ := NewMarkdown("Hello *world*")
		section, err := NewSection(text)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(section)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "section" {
			t.Errorf("got type %v, want section", result["type"])
		}

		textObj := result["text"].(map[string]any)
		if textObj["type"] != "mrkdwn" {
			t.Errorf("got text type %v, want mrkdwn", textObj["type"])
		}
	})

	t.Run("includes block_id when set", func(t *testing.T) {
		text, _ := NewMarkdown("Hello")
		section, _ := NewSection(text, WithSectionBlockID("section_1"))

		data, _ := json.Marshal(section)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["block_id"] != "section_1" {
			t.Errorf("got block_id %v, want 'section_1'", result["block_id"])
		}
	})

	t.Run("includes fields when set", func(t *testing.T) {
		text, _ := NewMarkdown("Main text")
		field1, _ := NewMarkdown("*Field 1*\nValue 1")
		field2, _ := NewMarkdown("*Field 2*\nValue 2")
		section, _ := NewSection(text, WithSectionFields(field1, field2))

		data, _ := json.Marshal(section)
		var result map[string]any
		mustUnmarshal(data, &result)

		fields := result["fields"].([]any)
		if len(fields) != 2 {
			t.Errorf("got %d fields, want 2", len(fields))
		}
	})

	t.Run("includes button accessory", func(t *testing.T) {
		text, _ := NewMarkdown("Click the button")
		btn, _ := NewButton("Click me")
		section, _ := NewSection(text, WithSectionAccessory(btn))

		data, _ := json.Marshal(section)
		var result map[string]any
		mustUnmarshal(data, &result)

		accessory := result["accessory"].(map[string]any)
		if accessory["type"] != "button" {
			t.Errorf("got accessory type %v, want 'button'", accessory["type"])
		}
	})

	t.Run("includes image accessory", func(t *testing.T) {
		text, _ := NewMarkdown("Check out this image")
		img, _ := NewImageElement("https://example.com/img.png", "An image")
		section, _ := NewSection(text, WithSectionAccessory(img))

		data, _ := json.Marshal(section)
		var result map[string]any
		mustUnmarshal(data, &result)

		accessory := result["accessory"].(map[string]any)
		if accessory["type"] != "image" {
			t.Errorf("got accessory type %v, want 'image'", accessory["type"])
		}
	})

	t.Run("includes select accessory", func(t *testing.T) {
		text, _ := NewMarkdown("Select an option")
		opt, _ := NewOption("Option", "opt")
		sel, _ := NewStaticSelect([]Option{opt})
		section, _ := NewSection(text, WithSectionAccessory(sel))

		data, _ := json.Marshal(section)
		var result map[string]any
		mustUnmarshal(data, &result)

		accessory := result["accessory"].(map[string]any)
		if accessory["type"] != "static_select" {
			t.Errorf("got accessory type %v, want 'static_select'", accessory["type"])
		}
	})

	t.Run("implements Block interface", func(t *testing.T) {
		text, _ := NewMarkdown("test")
		section, _ := NewSection(text)
		var _ Block = section
	})
}

func TestSectionWithFields(t *testing.T) {
	t.Run("creates valid section with fields only", func(t *testing.T) {
		field1, _ := NewMarkdown("*Name*\nJohn")
		field2, _ := NewMarkdown("*Age*\n30")
		section, err := NewSectionWithFields([]TextObject{field1, field2})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(section)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "section" {
			t.Errorf("got type %v, want section", result["type"])
		}

		if result["text"] != nil {
			t.Error("expected text to be nil for fields-only section")
		}

		fields := result["fields"].([]any)
		if len(fields) != 2 {
			t.Errorf("got %d fields, want 2", len(fields))
		}
	})

	t.Run("rejects empty fields", func(t *testing.T) {
		_, err := NewSectionWithFields([]TextObject{})
		if err == nil {
			t.Error("expected error for empty fields")
		}
	})
}
