package blocks

import (
	"encoding/json"
	"testing"
)

func TestButton(t *testing.T) {
	t.Run("creates valid button", func(t *testing.T) {
		btn, err := NewButton("Click me")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(btn)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "button" {
			t.Errorf("got type %v, want button", result["type"])
		}

		text := result["text"].(map[string]any)
		if text["text"] != "Click me" {
			t.Errorf("got text %v, want 'Click me'", text["text"])
		}
	})

	t.Run("includes action_id when set", func(t *testing.T) {
		btn, err := NewButton("Click me", WithActionID("btn_action"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, _ := json.Marshal(btn)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["action_id"] != "btn_action" {
			t.Errorf("got action_id %v, want 'btn_action'", result["action_id"])
		}
	})

	t.Run("includes value when set", func(t *testing.T) {
		btn, err := NewButton("Click me", WithValue("my_value"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, _ := json.Marshal(btn)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["value"] != "my_value" {
			t.Errorf("got value %v, want 'my_value'", result["value"])
		}
	})

	t.Run("includes style when set to primary", func(t *testing.T) {
		btn, err := NewButton("Click me", WithButtonStyle(ButtonStylePrimary))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, _ := json.Marshal(btn)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["style"] != "primary" {
			t.Errorf("got style %v, want 'primary'", result["style"])
		}
	})

	t.Run("includes style when set to danger", func(t *testing.T) {
		btn, err := NewButton("Delete", WithButtonStyle(ButtonStyleDanger))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, _ := json.Marshal(btn)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["style"] != "danger" {
			t.Errorf("got style %v, want 'danger'", result["style"])
		}
	})

	t.Run("includes URL when set", func(t *testing.T) {
		btn, err := NewButton("Visit", WithURL("https://example.com"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, _ := json.Marshal(btn)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["url"] != "https://example.com" {
			t.Errorf("got url %v, want 'https://example.com'", result["url"])
		}
	})

	t.Run("includes confirm dialog when set", func(t *testing.T) {
		confirm, _ := NewConfirmDialog("Confirm", "Are you sure?", "Yes", "No")
		btn, err := NewButton("Delete", WithButtonConfirm(confirm))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, _ := json.Marshal(btn)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["confirm"] == nil {
			t.Error("expected confirm dialog to be present")
		}
	})

	t.Run("rejects empty text", func(t *testing.T) {
		_, err := NewButton("")
		if err == nil {
			t.Error("expected error for empty text")
		}
	})

	t.Run("implements SectionAccessory interface", func(t *testing.T) {
		btn, _ := NewButton("test")
		var _ SectionAccessory = btn
	})

	t.Run("implements ActionsElement interface", func(t *testing.T) {
		btn, _ := NewButton("test")
		var _ ActionsElement = btn
	})
}
