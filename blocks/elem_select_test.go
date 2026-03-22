package blocks

import (
	"encoding/json"
	"testing"
)

func TestStaticSelect(t *testing.T) {
	t.Run("creates valid static select", func(t *testing.T) {
		opt1, _ := NewOption("Option 1", "opt1")
		opt2, _ := NewOption("Option 2", "opt2")

		sel, err := NewStaticSelect([]Option{opt1, opt2})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(sel)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "static_select" {
			t.Errorf("got type %v, want static_select", result["type"])
		}

		options := result["options"].([]any)
		if len(options) != 2 {
			t.Errorf("got %d options, want 2", len(options))
		}
	})

	t.Run("includes action_id when set", func(t *testing.T) {
		opt, _ := NewOption("Option", "opt")
		sel, _ := NewStaticSelect([]Option{opt}, WithStaticSelectActionID("select_action"))

		data, _ := json.Marshal(sel)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["action_id"] != "select_action" {
			t.Errorf("got action_id %v, want 'select_action'", result["action_id"])
		}
	})

	t.Run("includes placeholder when set", func(t *testing.T) {
		opt, _ := NewOption("Option", "opt")
		sel, _ := NewStaticSelect([]Option{opt}, WithStaticSelectPlaceholder("Choose..."))

		data, _ := json.Marshal(sel)
		var result map[string]any
		mustUnmarshal(data, &result)

		placeholder := result["placeholder"].(map[string]any)
		if placeholder["text"] != "Choose..." {
			t.Errorf("got placeholder %v, want 'Choose...'", placeholder["text"])
		}
	})

	t.Run("includes initial_option when set", func(t *testing.T) {
		opt1, _ := NewOption("Option 1", "opt1")
		opt2, _ := NewOption("Option 2", "opt2")
		sel, _ := NewStaticSelect([]Option{opt1, opt2}, WithStaticSelectInitialOption(opt1))

		data, _ := json.Marshal(sel)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["initial_option"] == nil {
			t.Error("expected initial_option to be present")
		}
	})

	t.Run("rejects empty options", func(t *testing.T) {
		_, err := NewStaticSelect([]Option{})
		if err == nil {
			t.Error("expected error for empty options")
		}
	})

	t.Run("implements SectionAccessory interface", func(t *testing.T) {
		opt, _ := NewOption("Option", "opt")
		sel, _ := NewStaticSelect([]Option{opt})
		var _ SectionAccessory = sel
	})

	t.Run("implements ActionsElement interface", func(t *testing.T) {
		opt, _ := NewOption("Option", "opt")
		sel, _ := NewStaticSelect([]Option{opt})
		var _ ActionsElement = sel
	})

	t.Run("implements InputElement interface", func(t *testing.T) {
		opt, _ := NewOption("Option", "opt")
		sel, _ := NewStaticSelect([]Option{opt})
		var _ InputElement = sel
	})
}

func TestMultiStaticSelect(t *testing.T) {
	t.Run("creates valid multi static select", func(t *testing.T) {
		opt1, _ := NewOption("Option 1", "opt1")
		opt2, _ := NewOption("Option 2", "opt2")

		sel, err := NewMultiStaticSelect([]Option{opt1, opt2})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(sel)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "multi_static_select" {
			t.Errorf("got type %v, want multi_static_select", result["type"])
		}
	})

	t.Run("includes initial_options when set", func(t *testing.T) {
		opt1, _ := NewOption("Option 1", "opt1")
		opt2, _ := NewOption("Option 2", "opt2")
		sel, _ := NewMultiStaticSelect([]Option{opt1, opt2},
			WithMultiStaticSelectInitialOptions(opt1))

		data, _ := json.Marshal(sel)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["initial_options"] == nil {
			t.Error("expected initial_options to be present")
		}
	})

	t.Run("includes max_selected_items when set", func(t *testing.T) {
		opt1, _ := NewOption("Option 1", "opt1")
		opt2, _ := NewOption("Option 2", "opt2")
		sel, _ := NewMultiStaticSelect([]Option{opt1, opt2},
			WithMultiStaticSelectMaxItems(3))

		data, _ := json.Marshal(sel)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["max_selected_items"] != float64(3) {
			t.Errorf("got max_selected_items %v, want 3", result["max_selected_items"])
		}
	})
}

func TestUsersSelect(t *testing.T) {
	t.Run("creates valid users select", func(t *testing.T) {
		sel := NewUsersSelect()

		data, err := json.Marshal(sel)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "users_select" {
			t.Errorf("got type %v, want users_select", result["type"])
		}
	})

	t.Run("includes initial_user when set", func(t *testing.T) {
		sel := NewUsersSelect(WithInitialUser("U12345"))

		data, _ := json.Marshal(sel)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["initial_user"] != "U12345" {
			t.Errorf("got initial_user %v, want 'U12345'", result["initial_user"])
		}
	})
}

func TestConversationsSelect(t *testing.T) {
	t.Run("creates valid conversations select", func(t *testing.T) {
		sel := NewConversationsSelect()

		data, err := json.Marshal(sel)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "conversations_select" {
			t.Errorf("got type %v, want conversations_select", result["type"])
		}
	})

	t.Run("includes default_to_current_conversation when set", func(t *testing.T) {
		sel := NewConversationsSelect(WithDefaultToCurrentConversation())

		data, _ := json.Marshal(sel)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["default_to_current_conversation"] != true {
			t.Error("expected default_to_current_conversation to be true")
		}
	})
}

func TestChannelsSelect(t *testing.T) {
	t.Run("creates valid channels select", func(t *testing.T) {
		sel := NewChannelsSelect()

		data, err := json.Marshal(sel)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "channels_select" {
			t.Errorf("got type %v, want channels_select", result["type"])
		}
	})

	t.Run("includes initial_channel when set", func(t *testing.T) {
		sel := NewChannelsSelect(WithInitialChannel("C12345"))

		data, _ := json.Marshal(sel)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["initial_channel"] != "C12345" {
			t.Errorf("got initial_channel %v, want 'C12345'", result["initial_channel"])
		}
	})
}
