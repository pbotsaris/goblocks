package blocks

import (
	"encoding/json"
	"testing"
)

func TestHeader(t *testing.T) {
	t.Run("creates valid header", func(t *testing.T) {
		header, err := NewHeader("My Header")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(header)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "header" {
			t.Errorf("got type %v, want header", result["type"])
		}

		text := result["text"].(map[string]any)
		if text["text"] != "My Header" {
			t.Errorf("got text %v, want 'My Header'", text["text"])
		}
		if text["type"] != "plain_text" {
			t.Errorf("got text type %v, want 'plain_text'", text["type"])
		}
	})

	t.Run("includes block_id when set", func(t *testing.T) {
		header, _ := NewHeader("Header", WithHeaderBlockID("header_1"))

		data, _ := json.Marshal(header)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["block_id"] != "header_1" {
			t.Errorf("got block_id %v, want 'header_1'", result["block_id"])
		}
	})

	t.Run("rejects empty text", func(t *testing.T) {
		_, err := NewHeader("")
		if err == nil {
			t.Error("expected error for empty text")
		}
	})

	t.Run("implements Block interface", func(t *testing.T) {
		header, _ := NewHeader("test")
		var _ Block = header
	})
}

func TestDivider(t *testing.T) {
	t.Run("creates valid divider", func(t *testing.T) {
		divider := NewDivider()

		data, err := json.Marshal(divider)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "divider" {
			t.Errorf("got type %v, want divider", result["type"])
		}
	})

	t.Run("includes block_id when set", func(t *testing.T) {
		divider := NewDivider(WithDividerBlockID("divider_1"))

		data, _ := json.Marshal(divider)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["block_id"] != "divider_1" {
			t.Errorf("got block_id %v, want 'divider_1'", result["block_id"])
		}
	})

	t.Run("implements Block interface", func(t *testing.T) {
		var _ Block = NewDivider()
	})
}

func TestImageBlock(t *testing.T) {
	t.Run("creates valid image block", func(t *testing.T) {
		img, err := NewImageBlock("https://example.com/image.png", "An image")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(img)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "image" {
			t.Errorf("got type %v, want image", result["type"])
		}
		if result["image_url"] != "https://example.com/image.png" {
			t.Errorf("got image_url %v, want 'https://example.com/image.png'", result["image_url"])
		}
		if result["alt_text"] != "An image" {
			t.Errorf("got alt_text %v, want 'An image'", result["alt_text"])
		}
	})

	t.Run("includes title when set", func(t *testing.T) {
		img, _ := NewImageBlock("https://example.com/image.png", "Alt",
			WithImageBlockTitle("My Image"))

		data, _ := json.Marshal(img)
		var result map[string]any
		mustUnmarshal(data, &result)

		title := result["title"].(map[string]any)
		if title["text"] != "My Image" {
			t.Errorf("got title %v, want 'My Image'", title["text"])
		}
	})

	t.Run("includes block_id when set", func(t *testing.T) {
		img, _ := NewImageBlock("https://example.com/image.png", "Alt",
			WithImageBlockID("image_1"))

		data, _ := json.Marshal(img)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["block_id"] != "image_1" {
			t.Errorf("got block_id %v, want 'image_1'", result["block_id"])
		}
	})

	t.Run("rejects empty URL", func(t *testing.T) {
		_, err := NewImageBlock("", "Alt text")
		if err == nil {
			t.Error("expected error for empty URL")
		}
	})

	t.Run("rejects empty alt_text", func(t *testing.T) {
		_, err := NewImageBlock("https://example.com/image.png", "")
		if err == nil {
			t.Error("expected error for empty alt_text")
		}
	})

	t.Run("implements Block interface", func(t *testing.T) {
		img, _ := NewImageBlock("https://example.com/image.png", "Alt")
		var _ Block = img
	})
}
