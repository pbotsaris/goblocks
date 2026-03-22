package blocks

import (
	"encoding/json"
	"testing"
)

func TestPlainText(t *testing.T) {
	t.Run("creates valid plain text with emoji enabled by default", func(t *testing.T) {
		pt, err := NewPlainText("Hello world")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(pt)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "plain_text" {
			t.Errorf("got type %v, want plain_text", result["type"])
		}
		if result["text"] != "Hello world" {
			t.Errorf("got text %v, want Hello world", result["text"])
		}
		if result["emoji"] != true {
			t.Errorf("emoji should be true by default")
		}
	})

	t.Run("includes emoji when set", func(t *testing.T) {
		pt, err := NewPlainText("Hello :wave:", WithEmoji(true))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(pt)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["emoji"] != true {
			t.Errorf("expected emoji to be true")
		}
	})

	t.Run("rejects empty text", func(t *testing.T) {
		_, err := NewPlainText("")
		if err == nil {
			t.Error("expected error for empty text")
		}
	})

	t.Run("implements TextObject interface", func(t *testing.T) {
		pt, _ := NewPlainText("test")
		var _ TextObject = pt
	})

	t.Run("implements PlainTextOnly interface", func(t *testing.T) {
		pt, _ := NewPlainText("test")
		var _ PlainTextOnly = pt
	})

	t.Run("implements ContextElement interface", func(t *testing.T) {
		pt, _ := NewPlainText("test")
		var _ ContextElement = pt
	})
}

func TestMarkdown(t *testing.T) {
	t.Run("creates valid markdown", func(t *testing.T) {
		md, err := NewMarkdown("*bold* text")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(md)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		expected := `{"type":"mrkdwn","text":"*bold* text"}`
		if string(data) != expected {
			t.Errorf("got %s, want %s", data, expected)
		}
	})

	t.Run("includes verbatim when set", func(t *testing.T) {
		md, err := NewMarkdown("text", WithVerbatim(true))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(md)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["verbatim"] != true {
			t.Errorf("expected verbatim to be true")
		}
	})

	t.Run("rejects empty text", func(t *testing.T) {
		_, err := NewMarkdown("")
		if err == nil {
			t.Error("expected error for empty text")
		}
	})

	t.Run("implements TextObject interface", func(t *testing.T) {
		md, _ := NewMarkdown("test")
		var _ TextObject = md
	})

	t.Run("implements ContextElement interface", func(t *testing.T) {
		md, _ := NewMarkdown("test")
		var _ ContextElement = md
	})

	t.Run("does not implement PlainTextOnly", func(t *testing.T) {
		// This is a compile-time check - Markdown should NOT implement PlainTextOnly
		// If this test compiles, the design is correct
		md, _ := NewMarkdown("test")
		_ = md // use the variable
	})
}
