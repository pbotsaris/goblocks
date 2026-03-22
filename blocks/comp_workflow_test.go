package blocks

import (
	"encoding/json"
	"testing"
)

func TestInputParameter(t *testing.T) {
	t.Run("creates valid input parameter", func(t *testing.T) {
		param, err := NewInputParameter("key", "value")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(param)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["name"] != "key" {
			t.Errorf("got name %v, want 'key'", result["name"])
		}
		if result["value"] != "value" {
			t.Errorf("got value %v, want 'value'", result["value"])
		}
	})

	t.Run("rejects empty name", func(t *testing.T) {
		_, err := NewInputParameter("", "value")
		if err == nil {
			t.Error("expected error for empty name")
		}
	})

	t.Run("rejects empty value", func(t *testing.T) {
		_, err := NewInputParameter("key", "")
		if err == nil {
			t.Error("expected error for empty value")
		}
	})
}

func TestTrigger(t *testing.T) {
	t.Run("creates valid trigger", func(t *testing.T) {
		trigger, err := NewTrigger("https://example.com/trigger")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(trigger)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["url"] != "https://example.com/trigger" {
			t.Errorf("got url %v, want 'https://example.com/trigger'", result["url"])
		}
	})

	t.Run("includes input parameters when set", func(t *testing.T) {
		param := MustInputParameter("key", "value")
		trigger, _ := NewTrigger("https://example.com/trigger",
			WithInputParameters(param))

		data, _ := json.Marshal(trigger)
		var result map[string]any
		mustUnmarshal(data, &result)

		params := result["customizable_input_parameters"].([]any)
		if len(params) != 1 {
			t.Errorf("got %d params, want 1", len(params))
		}
	})

	t.Run("rejects empty URL", func(t *testing.T) {
		_, err := NewTrigger("")
		if err == nil {
			t.Error("expected error for empty URL")
		}
	})
}

func TestWorkflow(t *testing.T) {
	t.Run("creates valid workflow", func(t *testing.T) {
		trigger := MustTrigger("https://example.com/trigger")
		workflow := NewWorkflow(trigger)

		data, err := json.Marshal(workflow)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		triggerObj := result["trigger"].(map[string]any)
		if triggerObj["url"] != "https://example.com/trigger" {
			t.Errorf("got trigger url %v, want 'https://example.com/trigger'", triggerObj["url"])
		}
	})
}
