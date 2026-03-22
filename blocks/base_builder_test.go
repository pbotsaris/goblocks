package blocks

import (
	"encoding/json"
	"testing"
)

func TestBuilder(t *testing.T) {
	t.Run("creates empty builder", func(t *testing.T) {
		b := NewBuilder()
		blocks, err := b.Build()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(blocks) != 0 {
			t.Errorf("got %d blocks, want 0", len(blocks))
		}
	})

	t.Run("adds section with text", func(t *testing.T) {
		text, _ := NewMarkdown("Hello *world*")
		b := NewBuilder().AddSection(text)

		blocks, err := b.Build()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(blocks) != 1 {
			t.Errorf("got %d blocks, want 1", len(blocks))
		}
	})

	t.Run("adds section with fields", func(t *testing.T) {
		field1, _ := NewMarkdown("*Name*\nJohn")
		field2, _ := NewMarkdown("*Age*\n30")
		b := NewBuilder().AddSectionWithFields([]TextObject{field1, field2})

		blocks, err := b.Build()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(blocks) != 1 {
			t.Errorf("got %d blocks, want 1", len(blocks))
		}
	})

	t.Run("adds divider", func(t *testing.T) {
		b := NewBuilder().AddDivider()

		blocks, err := b.Build()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(blocks) != 1 {
			t.Errorf("got %d blocks, want 1", len(blocks))
		}
	})

	t.Run("adds header", func(t *testing.T) {
		b := NewBuilder().AddHeader("My Header")

		blocks, err := b.Build()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(blocks) != 1 {
			t.Errorf("got %d blocks, want 1", len(blocks))
		}
	})

	t.Run("adds actions", func(t *testing.T) {
		btn, _ := NewButton("Click me")
		b := NewBuilder().AddActions([]ActionsElement{btn})

		blocks, err := b.Build()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(blocks) != 1 {
			t.Errorf("got %d blocks, want 1", len(blocks))
		}
	})

	t.Run("adds context", func(t *testing.T) {
		text, _ := NewMarkdown("Context info")
		b := NewBuilder().AddContext([]ContextElement{text})

		blocks, err := b.Build()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(blocks) != 1 {
			t.Errorf("got %d blocks, want 1", len(blocks))
		}
	})

	t.Run("adds input", func(t *testing.T) {
		pti := NewPlainTextInput()
		b := NewBuilder().AddInput("Enter text", pti)

		blocks, err := b.Build()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(blocks) != 1 {
			t.Errorf("got %d blocks, want 1", len(blocks))
		}
	})

	t.Run("adds image", func(t *testing.T) {
		b := NewBuilder().AddImage("https://example.com/img.png", "An image")

		blocks, err := b.Build()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(blocks) != 1 {
			t.Errorf("got %d blocks, want 1", len(blocks))
		}
	})

	t.Run("chains multiple blocks", func(t *testing.T) {
		text, _ := NewMarkdown("Hello")
		btn, _ := NewButton("Click")

		b := NewBuilder().
			AddHeader("Welcome").
			AddSection(text).
			AddDivider().
			AddActions([]ActionsElement{btn})

		blocks, err := b.Build()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(blocks) != 4 {
			t.Errorf("got %d blocks, want 4", len(blocks))
		}
	})

	t.Run("collects errors", func(t *testing.T) {
		b := NewBuilder().
			AddHeader(""). // empty header should error
			AddHeader("Valid")

		if !b.HasErrors() {
			t.Error("expected HasErrors() to be true")
		}

		errors := b.Errors()
		if len(errors) != 1 {
			t.Errorf("got %d errors, want 1", len(errors))
		}

		_, err := b.Build()
		if err == nil {
			t.Error("expected Build() to return error")
		}
	})

	t.Run("MustBuild panics on error", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected MustBuild to panic")
			}
		}()

		NewBuilder().
			AddHeader(""). // empty header should error
			MustBuild()
	})
}

func TestBuilderToModal(t *testing.T) {
	t.Run("converts to modal", func(t *testing.T) {
		text, _ := NewMarkdown("Modal content")
		b := NewBuilder().AddSection(text)

		modal, err := b.ToModal("My Modal")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, _ := json.Marshal(modal)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "modal" {
			t.Errorf("got type %v, want modal", result["type"])
		}

		title := result["title"].(map[string]any)
		if title["text"] != "My Modal" {
			t.Errorf("got title %v, want 'My Modal'", title["text"])
		}
	})

	t.Run("MustToModal panics on error", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected MustToModal to panic")
			}
		}()

		NewBuilder().
			AddHeader(""). // error
			MustToModal("Modal")
	})
}

func TestBuilderToMessage(t *testing.T) {
	t.Run("converts to message", func(t *testing.T) {
		text, _ := NewMarkdown("Message content")
		b := NewBuilder().AddSection(text)

		msg, err := b.ToMessage("Fallback text")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, _ := json.Marshal(msg)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["text"] != "Fallback text" {
			t.Errorf("got text %v, want 'Fallback text'", result["text"])
		}
	})

	t.Run("MustToMessage panics on error", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected MustToMessage to panic")
			}
		}()

		NewBuilder().
			AddHeader(""). // error
			MustToMessage("Fallback")
	})
}

func TestBuilderToHomeTab(t *testing.T) {
	t.Run("converts to home tab", func(t *testing.T) {
		text, _ := NewMarkdown("Home content")
		b := NewBuilder().AddSection(text)

		home, err := b.ToHomeTab()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, _ := json.Marshal(home)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "home" {
			t.Errorf("got type %v, want home", result["type"])
		}
	})

	t.Run("MustToHomeTab panics on error", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected MustToHomeTab to panic")
			}
		}()

		NewBuilder().
			AddHeader(""). // error
			MustToHomeTab()
	})
}

func TestBuilderJSON(t *testing.T) {
	t.Run("returns JSON with blocks wrapper", func(t *testing.T) {
		text, _ := NewMarkdown("Hello")
		b := NewBuilder().AddSection(text)

		data, err := b.JSON()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["blocks"] == nil {
			t.Error("expected 'blocks' key in JSON")
		}

		blocks := result["blocks"].([]any)
		if len(blocks) != 1 {
			t.Errorf("got %d blocks, want 1", len(blocks))
		}
	})

	t.Run("returns pretty JSON", func(t *testing.T) {
		text, _ := NewMarkdown("Hello")
		b := NewBuilder().AddSection(text)

		data, err := b.PrettyJSON()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Pretty JSON should have newlines and indentation
		if len(data) <= 50 {
			t.Error("expected pretty JSON to be formatted with indentation")
		}
	})

	t.Run("returns blocks-only JSON", func(t *testing.T) {
		text, _ := NewMarkdown("Hello")
		b := NewBuilder().AddSection(text)

		data, err := b.BlocksJSON()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var result []any
		mustUnmarshal(data, &result)

		if len(result) != 1 {
			t.Errorf("got %d blocks, want 1", len(result))
		}
	})
}

func TestBuilderAdd(t *testing.T) {
	t.Run("adds arbitrary block", func(t *testing.T) {
		divider := NewDivider()
		b := NewBuilder().Add(divider)

		blocks, err := b.Build()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(blocks) != 1 {
			t.Errorf("got %d blocks, want 1", len(blocks))
		}
	})
}
