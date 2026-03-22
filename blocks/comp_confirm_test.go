package blocks

import (
	"encoding/json"
	"testing"
)

func TestConfirmDialog(t *testing.T) {
	t.Run("creates valid confirm dialog", func(t *testing.T) {
		confirm, err := NewConfirmDialog("Confirm Action", "Are you sure?", "Yes", "No")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(confirm)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		title := result["title"].(map[string]any)
		if title["text"] != "Confirm Action" {
			t.Errorf("got title %v, want 'Confirm Action'", title["text"])
		}

		text := result["text"].(map[string]any)
		if text["text"] != "Are you sure?" {
			t.Errorf("got text %v, want 'Are you sure?'", text["text"])
		}

		confirm_btn := result["confirm"].(map[string]any)
		if confirm_btn["text"] != "Yes" {
			t.Errorf("got confirm %v, want 'Yes'", confirm_btn["text"])
		}

		deny := result["deny"].(map[string]any)
		if deny["text"] != "No" {
			t.Errorf("got deny %v, want 'No'", deny["text"])
		}
	})

	t.Run("includes style when set to danger", func(t *testing.T) {
		confirm, err := NewConfirmDialog("Delete", "Are you sure?", "Delete", "Cancel",
			WithConfirmStyle(ConfirmStyleDanger))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(confirm)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["style"] != "danger" {
			t.Errorf("got style %v, want 'danger'", result["style"])
		}
	})

	t.Run("rejects empty title", func(t *testing.T) {
		_, err := NewConfirmDialog("", "Text", "Confirm", "Deny")
		if err == nil {
			t.Error("expected error for empty title")
		}
	})

	t.Run("rejects empty text", func(t *testing.T) {
		_, err := NewConfirmDialog("Title", "", "Confirm", "Deny")
		if err == nil {
			t.Error("expected error for empty text")
		}
	})

	t.Run("rejects empty confirm", func(t *testing.T) {
		_, err := NewConfirmDialog("Title", "Text", "", "Deny")
		if err == nil {
			t.Error("expected error for empty confirm")
		}
	})

	t.Run("rejects empty deny", func(t *testing.T) {
		_, err := NewConfirmDialog("Title", "Text", "Confirm", "")
		if err == nil {
			t.Error("expected error for empty deny")
		}
	})
}
