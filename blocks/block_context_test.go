package blocks

import (
	"encoding/json"
	"testing"
)

func TestContext(t *testing.T) {
	t.Run("creates valid context block with text", func(t *testing.T) {
		text, _ := NewMarkdown("Context info")
		ctx, err := NewContext([]ContextElement{text})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(ctx)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "context" {
			t.Errorf("got type %v, want context", result["type"])
		}

		elements := result["elements"].([]any)
		if len(elements) != 1 {
			t.Errorf("got %d elements, want 1", len(elements))
		}
	})

	t.Run("includes block_id when set", func(t *testing.T) {
		text, _ := NewMarkdown("Info")
		ctx, _ := NewContext([]ContextElement{text}, WithContextBlockID("context_1"))

		data, _ := json.Marshal(ctx)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["block_id"] != "context_1" {
			t.Errorf("got block_id %v, want 'context_1'", result["block_id"])
		}
	})

	t.Run("supports image elements", func(t *testing.T) {
		img, _ := NewImageElement("https://example.com/icon.png", "Icon")
		text, _ := NewMarkdown("With icon")
		ctx, _ := NewContext([]ContextElement{img, text})

		data, _ := json.Marshal(ctx)
		var result map[string]any
		mustUnmarshal(data, &result)

		elements := result["elements"].([]any)
		if len(elements) != 2 {
			t.Errorf("got %d elements, want 2", len(elements))
		}

		elem0 := elements[0].(map[string]any)
		if elem0["type"] != "image" {
			t.Errorf("element 0 type = %v, want image", elem0["type"])
		}
	})

	t.Run("supports PlainText elements", func(t *testing.T) {
		text, _ := NewPlainText("Plain context")
		ctx, _ := NewContext([]ContextElement{text})

		data, _ := json.Marshal(ctx)
		var result map[string]any
		mustUnmarshal(data, &result)

		elements := result["elements"].([]any)
		elem := elements[0].(map[string]any)
		if elem["type"] != "plain_text" {
			t.Errorf("element type = %v, want plain_text", elem["type"])
		}
	})

	t.Run("rejects empty elements", func(t *testing.T) {
		_, err := NewContext([]ContextElement{})
		if err == nil {
			t.Error("expected error for empty elements")
		}
	})

	t.Run("implements Block interface", func(t *testing.T) {
		text, _ := NewMarkdown("test")
		ctx, _ := NewContext([]ContextElement{text})
		var _ Block = ctx
	})
}
