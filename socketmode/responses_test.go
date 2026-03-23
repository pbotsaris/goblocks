package socketmode

import (
	"encoding/json"
	"testing"

	"github.com/pbotsaris/goblocks/blocks"
)

func TestEmptyResponse(t *testing.T) {
	t.Run("toPayload returns nil", func(t *testing.T) {
		resp := EmptyResponse{}
		if resp.toPayload() != nil {
			t.Error("expected nil payload")
		}
	})

	t.Run("implements Response interface", func(t *testing.T) {
		var _ Response = EmptyResponse{}
	})
}

func TestMessageResponse(t *testing.T) {
	t.Run("toPayload returns message", func(t *testing.T) {
		text, _ := blocks.NewMarkdown("Hello *world*")
		section, _ := blocks.NewSection(text)
		msg, _ := blocks.NewMessage("fallback", []blocks.Block{section})

		resp := MessageResponse{message: msg}
		payload := resp.toPayload()

		// Marshal to JSON to verify structure
		data, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(t, data, &result)

		blks := result["blocks"].([]any)
		if len(blks) != 1 {
			t.Errorf("got %d blocks, want 1", len(blks))
		}
	})

	t.Run("implements Response interface", func(t *testing.T) {
		var _ Response = MessageResponse{}
	})
}

func TestModalResponse(t *testing.T) {
	t.Run("update action", func(t *testing.T) {
		text, _ := blocks.NewMarkdown("Content")
		section, _ := blocks.NewSection(text)
		modal, _ := blocks.NewModal("Title", []blocks.Block{section})

		resp := ModalResponse{
			action: ModalActionUpdate,
			view:   &modal,
		}
		payload := resp.toPayload()

		data, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(t, data, &result)

		if result["response_action"] != "update" {
			t.Errorf("got response_action %v, want 'update'", result["response_action"])
		}
		if result["view"] == nil {
			t.Error("expected view to be present")
		}
	})

	t.Run("push action", func(t *testing.T) {
		text, _ := blocks.NewMarkdown("Content")
		section, _ := blocks.NewSection(text)
		modal, _ := blocks.NewModal("Title", []blocks.Block{section})

		resp := ModalResponse{
			action: ModalActionPush,
			view:   &modal,
		}
		payload := resp.toPayload()

		data, _ := json.Marshal(payload)
		var result map[string]any
		mustUnmarshal(t, data, &result)

		if result["response_action"] != "push" {
			t.Errorf("got response_action %v, want 'push'", result["response_action"])
		}
	})

	t.Run("clear action", func(t *testing.T) {
		resp := ModalResponse{
			action: ModalActionClear,
		}
		payload := resp.toPayload()

		data, _ := json.Marshal(payload)
		var result map[string]any
		mustUnmarshal(t, data, &result)

		if result["response_action"] != "clear" {
			t.Errorf("got response_action %v, want 'clear'", result["response_action"])
		}
		if result["view"] != nil {
			t.Error("expected view to be absent for clear action")
		}
	})

	t.Run("errors action", func(t *testing.T) {
		resp := ModalResponse{
			action: ModalActionErrors,
			errors: map[string]string{
				"email_block": "Invalid email address",
				"name_block":  "Name is required",
			},
		}
		payload := resp.toPayload()

		data, _ := json.Marshal(payload)
		var result map[string]any
		mustUnmarshal(t, data, &result)

		if result["response_action"] != "errors" {
			t.Errorf("got response_action %v, want 'errors'", result["response_action"])
		}

		errors := result["errors"].(map[string]any)
		if errors["email_block"] != "Invalid email address" {
			t.Errorf("got error %v", errors["email_block"])
		}
	})

	t.Run("nil view returns nil for update", func(t *testing.T) {
		resp := ModalResponse{
			action: ModalActionUpdate,
			view:   nil,
		}
		if resp.toPayload() != nil {
			t.Error("expected nil payload for nil view")
		}
	})

	t.Run("empty errors returns nil", func(t *testing.T) {
		resp := ModalResponse{
			action: ModalActionErrors,
			errors: nil,
		}
		if resp.toPayload() != nil {
			t.Error("expected nil payload for nil errors")
		}
	})

	t.Run("implements Response interface", func(t *testing.T) {
		var _ Response = ModalResponse{}
	})
}

func TestOptionsResponse(t *testing.T) {
	t.Run("options payload", func(t *testing.T) {
		opt1, _ := blocks.NewOption("Option 1", "opt1")
		opt2, _ := blocks.NewOption("Option 2", "opt2")

		resp := OptionsResponse{
			options: []blocks.Option{opt1, opt2},
		}
		payload := resp.toPayload()

		data, _ := json.Marshal(payload)
		var result map[string]any
		mustUnmarshal(t, data, &result)

		opts := result["options"].([]any)
		if len(opts) != 2 {
			t.Errorf("got %d options, want 2", len(opts))
		}
	})

	t.Run("option_groups payload", func(t *testing.T) {
		opt1, _ := blocks.NewOption("Option 1", "opt1")
		opt2, _ := blocks.NewOption("Option 2", "opt2")
		group, _ := blocks.NewOptionGroup("Group 1", []blocks.Option{opt1, opt2})

		resp := OptionsResponse{
			optionGroups: []blocks.OptionGroup{group},
		}
		payload := resp.toPayload()

		data, _ := json.Marshal(payload)
		var result map[string]any
		mustUnmarshal(t, data, &result)

		groups := result["option_groups"].([]any)
		if len(groups) != 1 {
			t.Errorf("got %d groups, want 1", len(groups))
		}
	})

	t.Run("option_groups takes precedence", func(t *testing.T) {
		opt1, _ := blocks.NewOption("Option 1", "opt1")
		group, _ := blocks.NewOptionGroup("Group 1", []blocks.Option{opt1})

		resp := OptionsResponse{
			options:      []blocks.Option{opt1},
			optionGroups: []blocks.OptionGroup{group},
		}
		payload := resp.toPayload()

		data, _ := json.Marshal(payload)
		var result map[string]any
		mustUnmarshal(t, data, &result)

		if result["option_groups"] == nil {
			t.Error("expected option_groups to be present")
		}
		if result["options"] != nil {
			t.Error("expected options to be absent when option_groups present")
		}
	})

	t.Run("implements Response interface", func(t *testing.T) {
		var _ Response = OptionsResponse{}
	})
}

func TestResponseBuilders(t *testing.T) {
	t.Run("NoResponse returns EmptyResponse", func(t *testing.T) {
		resp := NoResponse()
		if _, ok := resp.(EmptyResponse); !ok {
			t.Error("expected EmptyResponse")
		}
	})

	t.Run("RespondWithMessage returns MessageResponse", func(t *testing.T) {
		msg, _ := blocks.NewMessage("test", []blocks.Block{})
		resp := RespondWithMessage(msg)
		if _, ok := resp.(MessageResponse); !ok {
			t.Error("expected MessageResponse")
		}
	})

	t.Run("RespondWithBlocks returns MessageResponse", func(t *testing.T) {
		text, _ := blocks.NewMarkdown("Content")
		section, _ := blocks.NewSection(text)
		resp := RespondWithBlocks([]blocks.Block{section})
		if _, ok := resp.(MessageResponse); !ok {
			t.Error("expected MessageResponse")
		}
	})

	t.Run("RespondWithBlocks handles invalid blocks", func(t *testing.T) {
		// Empty blocks should work
		resp := RespondWithBlocks([]blocks.Block{})
		if _, ok := resp.(MessageResponse); !ok {
			t.Error("expected MessageResponse for empty blocks")
		}
	})

	t.Run("RespondWithModalUpdate returns correct type", func(t *testing.T) {
		modal, _ := blocks.NewModal("Title", []blocks.Block{})
		resp := RespondWithModalUpdate(modal)
		mr, ok := resp.(ModalResponse)
		if !ok {
			t.Fatal("expected ModalResponse")
		}
		if mr.action != ModalActionUpdate {
			t.Error("expected ModalActionUpdate")
		}
	})

	t.Run("RespondWithModalPush returns correct type", func(t *testing.T) {
		modal, _ := blocks.NewModal("Title", []blocks.Block{})
		resp := RespondWithModalPush(modal)
		mr, ok := resp.(ModalResponse)
		if !ok {
			t.Fatal("expected ModalResponse")
		}
		if mr.action != ModalActionPush {
			t.Error("expected ModalActionPush")
		}
	})

	t.Run("RespondWithModalClear returns correct type", func(t *testing.T) {
		resp := RespondWithModalClear()
		mr, ok := resp.(ModalResponse)
		if !ok {
			t.Fatal("expected ModalResponse")
		}
		if mr.action != ModalActionClear {
			t.Error("expected ModalActionClear")
		}
	})

	t.Run("RespondWithErrors returns correct type", func(t *testing.T) {
		errors := map[string]string{"field": "error"}
		resp := RespondWithErrors(errors)
		mr, ok := resp.(ModalResponse)
		if !ok {
			t.Fatal("expected ModalResponse")
		}
		if mr.action != ModalActionErrors {
			t.Error("expected ModalActionErrors")
		}
		if mr.errors["field"] != "error" {
			t.Error("expected errors to be set")
		}
	})

	t.Run("RespondWithOptions returns correct type", func(t *testing.T) {
		opt, _ := blocks.NewOption("Option", "value")
		resp := RespondWithOptions([]blocks.Option{opt})
		or, ok := resp.(OptionsResponse)
		if !ok {
			t.Fatal("expected OptionsResponse")
		}
		if len(or.options) != 1 {
			t.Error("expected 1 option")
		}
	})

	t.Run("RespondWithOptionGroups returns correct type", func(t *testing.T) {
		opt, _ := blocks.NewOption("Option", "value")
		group, _ := blocks.NewOptionGroup("Group", []blocks.Option{opt})
		resp := RespondWithOptionGroups([]blocks.OptionGroup{group})
		or, ok := resp.(OptionsResponse)
		if !ok {
			t.Fatal("expected OptionsResponse")
		}
		if len(or.optionGroups) != 1 {
			t.Error("expected 1 option group")
		}
	})
}
