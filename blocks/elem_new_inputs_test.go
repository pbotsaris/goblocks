package blocks

import (
	"encoding/json"
	"testing"
)

func TestEmailInput(t *testing.T) {
	t.Run("creates valid email input", func(t *testing.T) {
		email := NewEmailInput()

		data, err := json.Marshal(email)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "email_text_input" {
			t.Errorf("got type %v, want 'email_text_input'", result["type"])
		}
	})

	t.Run("includes all options when set", func(t *testing.T) {
		email := NewEmailInput(
			WithEmailInputActionID("email_action"),
			WithEmailInputInitialValue("test@example.com"),
			WithEmailInputFocusOnLoad(),
			WithEmailInputPlaceholder("Enter email"),
		)

		data, _ := json.Marshal(email)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["action_id"] != "email_action" {
			t.Errorf("got action_id %v, want 'email_action'", result["action_id"])
		}
		if result["initial_value"] != "test@example.com" {
			t.Errorf("got initial_value %v, want 'test@example.com'", result["initial_value"])
		}
		if result["focus_on_load"] != true {
			t.Errorf("got focus_on_load %v, want true", result["focus_on_load"])
		}
	})

	t.Run("implements InputElement interface", func(t *testing.T) {
		var _ InputElement = NewEmailInput()
	})
}

func TestNumberInput(t *testing.T) {
	t.Run("creates valid number input", func(t *testing.T) {
		num := NewNumberInput(true)

		data, err := json.Marshal(num)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "number_input" {
			t.Errorf("got type %v, want 'number_input'", result["type"])
		}
		if result["is_decimal_allowed"] != true {
			t.Errorf("got is_decimal_allowed %v, want true", result["is_decimal_allowed"])
		}
	})

	t.Run("includes min/max values when set", func(t *testing.T) {
		num := NewNumberInput(false,
			WithNumberInputMinValue("0"),
			WithNumberInputMaxValue("100"),
			WithNumberInputInitialValue("50"),
		)

		data, _ := json.Marshal(num)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["min_value"] != "0" {
			t.Errorf("got min_value %v, want '0'", result["min_value"])
		}
		if result["max_value"] != "100" {
			t.Errorf("got max_value %v, want '100'", result["max_value"])
		}
		if result["initial_value"] != "50" {
			t.Errorf("got initial_value %v, want '50'", result["initial_value"])
		}
	})

	t.Run("implements InputElement interface", func(t *testing.T) {
		var _ InputElement = NewNumberInput(true)
	})
}

func TestURLInput(t *testing.T) {
	t.Run("creates valid URL input", func(t *testing.T) {
		urlInput := NewURLInput()

		data, err := json.Marshal(urlInput)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "url_text_input" {
			t.Errorf("got type %v, want 'url_text_input'", result["type"])
		}
	})

	t.Run("includes all options when set", func(t *testing.T) {
		urlInput := NewURLInput(
			WithURLInputActionID("url_action"),
			WithURLInputInitialValue("https://example.com"),
			WithURLInputFocusOnLoad(),
		)

		data, _ := json.Marshal(urlInput)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["action_id"] != "url_action" {
			t.Errorf("got action_id %v, want 'url_action'", result["action_id"])
		}
		if result["initial_value"] != "https://example.com" {
			t.Errorf("got initial_value %v, want 'https://example.com'", result["initial_value"])
		}
	})

	t.Run("implements InputElement interface", func(t *testing.T) {
		var _ InputElement = NewURLInput()
	})
}

func TestFileInput(t *testing.T) {
	t.Run("creates valid file input", func(t *testing.T) {
		fileInput := NewFileInput()

		data, err := json.Marshal(fileInput)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "file_input" {
			t.Errorf("got type %v, want 'file_input'", result["type"])
		}
	})

	t.Run("includes filetypes and max_files when set", func(t *testing.T) {
		fileInput := NewFileInput(
			WithFileInputActionID("file_action"),
			WithFileInputFiletypes([]string{"pdf", "doc"}),
			WithFileInputMaxFiles(5),
		)

		data, _ := json.Marshal(fileInput)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["action_id"] != "file_action" {
			t.Errorf("got action_id %v, want 'file_action'", result["action_id"])
		}
		filetypes := result["filetypes"].([]any)
		if len(filetypes) != 2 {
			t.Errorf("got %d filetypes, want 2", len(filetypes))
		}
		if result["max_files"] != float64(5) {
			t.Errorf("got max_files %v, want 5", result["max_files"])
		}
	})

	t.Run("implements InputElement interface", func(t *testing.T) {
		var _ InputElement = NewFileInput()
	})
}

func TestRichTextInput(t *testing.T) {
	t.Run("creates valid rich text input", func(t *testing.T) {
		rtInput := NewRichTextInput()

		data, err := json.Marshal(rtInput)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "rich_text_input" {
			t.Errorf("got type %v, want 'rich_text_input'", result["type"])
		}
	})

	t.Run("includes all options when set", func(t *testing.T) {
		rtInput := NewRichTextInput(
			WithRichTextInputActionID("rt_action"),
			WithRichTextInputFocusOnLoad(),
			WithRichTextInputPlaceholder("Enter text"),
		)

		data, _ := json.Marshal(rtInput)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["action_id"] != "rt_action" {
			t.Errorf("got action_id %v, want 'rt_action'", result["action_id"])
		}
		if result["focus_on_load"] != true {
			t.Errorf("got focus_on_load %v, want true", result["focus_on_load"])
		}
	})

	t.Run("implements InputElement interface", func(t *testing.T) {
		var _ InputElement = NewRichTextInput()
	})
}

func TestWorkflowButton(t *testing.T) {
	t.Run("creates valid workflow button", func(t *testing.T) {
		trigger := MustTrigger("https://example.com/trigger")
		workflow := NewWorkflow(trigger)
		btn, err := NewWorkflowButton("Start Workflow", workflow)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(btn)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "workflow_button" {
			t.Errorf("got type %v, want 'workflow_button'", result["type"])
		}
		text := result["text"].(map[string]any)
		if text["text"] != "Start Workflow" {
			t.Errorf("got text %v, want 'Start Workflow'", text["text"])
		}
	})

	t.Run("includes style when set", func(t *testing.T) {
		trigger := MustTrigger("https://example.com/trigger")
		workflow := NewWorkflow(trigger)
		btn, _ := NewWorkflowButton("Start", workflow,
			WithWorkflowButtonStyle(ButtonStylePrimary),
			WithWorkflowButtonActionID("wf_btn"),
		)

		data, _ := json.Marshal(btn)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["style"] != "primary" {
			t.Errorf("got style %v, want 'primary'", result["style"])
		}
		if result["action_id"] != "wf_btn" {
			t.Errorf("got action_id %v, want 'wf_btn'", result["action_id"])
		}
	})

	t.Run("rejects empty text", func(t *testing.T) {
		trigger := MustTrigger("https://example.com/trigger")
		workflow := NewWorkflow(trigger)
		_, err := NewWorkflowButton("", workflow)
		if err == nil {
			t.Error("expected error for empty text")
		}
	})

	t.Run("implements SectionAccessory interface", func(t *testing.T) {
		trigger := MustTrigger("https://example.com/trigger")
		workflow := NewWorkflow(trigger)
		btn, _ := NewWorkflowButton("Test", workflow)
		var _ SectionAccessory = btn
	})

	t.Run("implements ActionsElement interface", func(t *testing.T) {
		trigger := MustTrigger("https://example.com/trigger")
		workflow := NewWorkflow(trigger)
		btn, _ := NewWorkflowButton("Test", workflow)
		var _ ActionsElement = btn
	})
}

func TestFeedbackButtons(t *testing.T) {
	t.Run("creates valid feedback buttons", func(t *testing.T) {
		fb := NewFeedbackButtons()

		data, err := json.Marshal(fb)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "feedback_buttons" {
			t.Errorf("got type %v, want 'feedback_buttons'", result["type"])
		}
	})

	t.Run("includes action_id when set", func(t *testing.T) {
		fb := NewFeedbackButtons(WithFeedbackButtonsActionID("fb_action"))

		data, _ := json.Marshal(fb)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["action_id"] != "fb_action" {
			t.Errorf("got action_id %v, want 'fb_action'", result["action_id"])
		}
	})

	t.Run("implements ContextActionsElement interface", func(t *testing.T) {
		var _ ContextActionsElement = NewFeedbackButtons()
	})
}

func TestIconButton(t *testing.T) {
	t.Run("creates valid icon button", func(t *testing.T) {
		icon := NewIcon("copy")
		btn := NewIconButton(icon)

		data, err := json.Marshal(btn)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "icon_button" {
			t.Errorf("got type %v, want 'icon_button'", result["type"])
		}
		iconObj := result["icon"].(map[string]any)
		if iconObj["name"] != "copy" {
			t.Errorf("got icon name %v, want 'copy'", iconObj["name"])
		}
	})

	t.Run("includes all options when set", func(t *testing.T) {
		icon := NewIcon("copy")
		btn := NewIconButton(icon,
			WithIconButtonActionID("icon_action"),
			WithIconButtonAltText("Copy text"),
		)

		data, _ := json.Marshal(btn)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["action_id"] != "icon_action" {
			t.Errorf("got action_id %v, want 'icon_action'", result["action_id"])
		}
		if result["alt_text"] != "Copy text" {
			t.Errorf("got alt_text %v, want 'Copy text'", result["alt_text"])
		}
	})

	t.Run("implements ContextActionsElement interface", func(t *testing.T) {
		icon := NewIcon("copy")
		var _ ContextActionsElement = NewIconButton(icon)
	})
}
