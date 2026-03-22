package blocks

import (
	"encoding/json"
	"testing"
)

func TestModal(t *testing.T) {
	t.Run("creates valid modal with blocks", func(t *testing.T) {
		text, _ := NewMarkdown("Modal content")
		section, _ := NewSection(text)
		blocks := []Block{section}

		modal, err := NewModal("My Modal", blocks)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(modal)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "modal" {
			t.Errorf("got type %v, want modal", result["type"])
		}

		title := result["title"].(map[string]any)
		if title["text"] != "My Modal" {
			t.Errorf("got title %v, want 'My Modal'", title["text"])
		}

		blks := result["blocks"].([]any)
		if len(blks) != 1 {
			t.Errorf("got %d blocks, want 1", len(blks))
		}
	})

	t.Run("includes submit button when set", func(t *testing.T) {
		text, _ := NewMarkdown("Content")
		section, _ := NewSection(text)
		modal, _ := NewModal("Modal", []Block{section}, WithModalSubmit("Submit"))

		data, _ := json.Marshal(modal)
		var result map[string]any
		mustUnmarshal(data, &result)

		submit := result["submit"].(map[string]any)
		if submit["text"] != "Submit" {
			t.Errorf("got submit %v, want 'Submit'", submit["text"])
		}
	})

	t.Run("includes close button when set", func(t *testing.T) {
		text, _ := NewMarkdown("Content")
		section, _ := NewSection(text)
		modal, _ := NewModal("Modal", []Block{section}, WithModalClose("Cancel"))

		data, _ := json.Marshal(modal)
		var result map[string]any
		mustUnmarshal(data, &result)

		close := result["close"].(map[string]any)
		if close["text"] != "Cancel" {
			t.Errorf("got close %v, want 'Cancel'", close["text"])
		}
	})

	t.Run("includes private_metadata when set", func(t *testing.T) {
		text, _ := NewMarkdown("Content")
		section, _ := NewSection(text)
		modal, _ := NewModal("Modal", []Block{section},
			WithModalPrivateMetadata(`{"key":"value"}`))

		data, _ := json.Marshal(modal)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["private_metadata"] != `{"key":"value"}` {
			t.Errorf("got private_metadata %v, want '{\"key\":\"value\"}'", result["private_metadata"])
		}
	})

	t.Run("includes callback_id when set", func(t *testing.T) {
		text, _ := NewMarkdown("Content")
		section, _ := NewSection(text)
		modal, _ := NewModal("Modal", []Block{section},
			WithModalCallbackID("modal_callback"))

		data, _ := json.Marshal(modal)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["callback_id"] != "modal_callback" {
			t.Errorf("got callback_id %v, want 'modal_callback'", result["callback_id"])
		}
	})

	t.Run("includes clear_on_close when set", func(t *testing.T) {
		text, _ := NewMarkdown("Content")
		section, _ := NewSection(text)
		modal, _ := NewModal("Modal", []Block{section}, WithModalClearOnClose())

		data, _ := json.Marshal(modal)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["clear_on_close"] != true {
			t.Error("expected clear_on_close to be true")
		}
	})

	t.Run("includes notify_on_close when set", func(t *testing.T) {
		text, _ := NewMarkdown("Content")
		section, _ := NewSection(text)
		modal, _ := NewModal("Modal", []Block{section}, WithModalNotifyOnClose())

		data, _ := json.Marshal(modal)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["notify_on_close"] != true {
			t.Error("expected notify_on_close to be true")
		}
	})

	t.Run("includes submit_disabled when set", func(t *testing.T) {
		text, _ := NewMarkdown("Content")
		section, _ := NewSection(text)
		modal, _ := NewModal("Modal", []Block{section}, WithModalSubmitDisabled())

		data, _ := json.Marshal(modal)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["submit_disabled"] != true {
			t.Error("expected submit_disabled to be true")
		}
	})

	t.Run("rejects empty title", func(t *testing.T) {
		text, _ := NewMarkdown("Content")
		section, _ := NewSection(text)
		_, err := NewModal("", []Block{section})
		if err == nil {
			t.Error("expected error for empty title")
		}
	})

	t.Run("allows empty blocks", func(t *testing.T) {
		// Slack API allows modals with no blocks
		modal, err := NewModal("Modal", []Block{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, _ := json.Marshal(modal)
		var result map[string]any
		mustUnmarshal(data, &result)

		blks := result["blocks"].([]any)
		if len(blks) != 0 {
			t.Errorf("got %d blocks, want 0", len(blks))
		}
	})
}

func TestMessage(t *testing.T) {
	t.Run("creates valid message with blocks", func(t *testing.T) {
		text, _ := NewMarkdown("Hello")
		section, _ := NewSection(text)
		blocks := []Block{section}

		msg, err := NewMessage("Fallback text", blocks)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(msg)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["text"] != "Fallback text" {
			t.Errorf("got text %v, want 'Fallback text'", result["text"])
		}

		blks := result["blocks"].([]any)
		if len(blks) != 1 {
			t.Errorf("got %d blocks, want 1", len(blks))
		}
	})

	t.Run("includes thread_ts when set", func(t *testing.T) {
		text, _ := NewMarkdown("Reply")
		section, _ := NewSection(text)
		msg, _ := NewMessage("Fallback", []Block{section},
			WithMessageThreadTS("1234567890.123456"))

		data, _ := json.Marshal(msg)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["thread_ts"] != "1234567890.123456" {
			t.Errorf("got thread_ts %v, want '1234567890.123456'", result["thread_ts"])
		}
	})

	t.Run("includes mrkdwn when set", func(t *testing.T) {
		text, _ := NewMarkdown("*Bold*")
		section, _ := NewSection(text)
		msg, _ := NewMessage("Fallback", []Block{section}, WithMessageMrkdwn())

		data, _ := json.Marshal(msg)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["mrkdwn"] != true {
			t.Error("expected mrkdwn to be true")
		}
	})

	t.Run("allows empty fallback text", func(t *testing.T) {
		// Slack API allows messages with blocks but no text
		text, _ := NewMarkdown("Content")
		section, _ := NewSection(text)
		msg, err := NewMessage("", []Block{section})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, _ := json.Marshal(msg)
		var result map[string]any
		mustUnmarshal(data, &result)

		// Empty text is omitted from JSON (omitempty)
		if result["text"] != nil && result["text"] != "" {
			t.Errorf("expected text to be nil or empty, got %v", result["text"])
		}
	})
}

func TestHomeTab(t *testing.T) {
	t.Run("creates valid home tab with blocks", func(t *testing.T) {
		text, _ := NewMarkdown("Welcome home!")
		section, _ := NewSection(text)
		blocks := []Block{section}

		home, err := NewHomeTab(blocks)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(home)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "home" {
			t.Errorf("got type %v, want home", result["type"])
		}

		blks := result["blocks"].([]any)
		if len(blks) != 1 {
			t.Errorf("got %d blocks, want 1", len(blks))
		}
	})

	t.Run("includes private_metadata when set", func(t *testing.T) {
		text, _ := NewMarkdown("Content")
		section, _ := NewSection(text)
		home, _ := NewHomeTab([]Block{section},
			WithHomeTabPrivateMetadata(`{"user":"123"}`))

		data, _ := json.Marshal(home)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["private_metadata"] != `{"user":"123"}` {
			t.Errorf("got private_metadata %v", result["private_metadata"])
		}
	})

	t.Run("includes callback_id when set", func(t *testing.T) {
		text, _ := NewMarkdown("Content")
		section, _ := NewSection(text)
		home, _ := NewHomeTab([]Block{section},
			WithHomeTabCallbackID("home_callback"))

		data, _ := json.Marshal(home)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["callback_id"] != "home_callback" {
			t.Errorf("got callback_id %v, want 'home_callback'", result["callback_id"])
		}
	})

	t.Run("includes external_id when set", func(t *testing.T) {
		text, _ := NewMarkdown("Content")
		section, _ := NewSection(text)
		home, _ := NewHomeTab([]Block{section},
			WithHomeTabExternalID("ext_123"))

		data, _ := json.Marshal(home)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["external_id"] != "ext_123" {
			t.Errorf("got external_id %v, want 'ext_123'", result["external_id"])
		}
	})

	t.Run("allows empty blocks", func(t *testing.T) {
		// Slack API allows home tabs with no blocks
		home, err := NewHomeTab([]Block{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, _ := json.Marshal(home)
		var result map[string]any
		mustUnmarshal(data, &result)

		blks := result["blocks"].([]any)
		if len(blks) != 0 {
			t.Errorf("got %d blocks, want 0", len(blks))
		}
	})
}
