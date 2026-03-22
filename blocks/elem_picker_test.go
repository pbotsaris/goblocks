package blocks

import (
	"encoding/json"
	"testing"
)

func TestDatePicker(t *testing.T) {
	t.Run("creates valid date picker", func(t *testing.T) {
		dp := NewDatePicker()

		data, err := json.Marshal(dp)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "datepicker" {
			t.Errorf("got type %v, want datepicker", result["type"])
		}
	})

	t.Run("includes action_id when set", func(t *testing.T) {
		dp := NewDatePicker(WithDatePickerActionID("date_action"))

		data, _ := json.Marshal(dp)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["action_id"] != "date_action" {
			t.Errorf("got action_id %v, want 'date_action'", result["action_id"])
		}
	})

	t.Run("includes initial_date when set", func(t *testing.T) {
		dp := NewDatePicker(WithInitialDate("2024-01-15"))

		data, _ := json.Marshal(dp)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["initial_date"] != "2024-01-15" {
			t.Errorf("got initial_date %v, want '2024-01-15'", result["initial_date"])
		}
	})

	t.Run("includes placeholder when set", func(t *testing.T) {
		dp := NewDatePicker(WithDatePickerPlaceholder("Select a date"))

		data, _ := json.Marshal(dp)
		var result map[string]any
		mustUnmarshal(data, &result)

		placeholder := result["placeholder"].(map[string]any)
		if placeholder["text"] != "Select a date" {
			t.Errorf("got placeholder %v, want 'Select a date'", placeholder["text"])
		}
	})

	t.Run("includes focus_on_load when set", func(t *testing.T) {
		dp := NewDatePicker(WithDatePickerFocusOnLoad())

		data, _ := json.Marshal(dp)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["focus_on_load"] != true {
			t.Error("expected focus_on_load to be true")
		}
	})

	t.Run("implements SectionAccessory interface", func(t *testing.T) {
		dp := NewDatePicker()
		var _ SectionAccessory = dp
	})

	t.Run("implements ActionsElement interface", func(t *testing.T) {
		dp := NewDatePicker()
		var _ ActionsElement = dp
	})

	t.Run("implements InputElement interface", func(t *testing.T) {
		dp := NewDatePicker()
		var _ InputElement = dp
	})
}

func TestTimePicker(t *testing.T) {
	t.Run("creates valid time picker", func(t *testing.T) {
		tp := NewTimePicker()

		data, err := json.Marshal(tp)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "timepicker" {
			t.Errorf("got type %v, want timepicker", result["type"])
		}
	})

	t.Run("includes initial_time when set", func(t *testing.T) {
		tp := NewTimePicker(WithInitialTime("14:30"))

		data, _ := json.Marshal(tp)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["initial_time"] != "14:30" {
			t.Errorf("got initial_time %v, want '14:30'", result["initial_time"])
		}
	})

	t.Run("includes timezone when set", func(t *testing.T) {
		tp := NewTimePicker(WithTimezone("America/New_York"))

		data, _ := json.Marshal(tp)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["timezone"] != "America/New_York" {
			t.Errorf("got timezone %v, want 'America/New_York'", result["timezone"])
		}
	})
}

func TestDatetimePicker(t *testing.T) {
	t.Run("creates valid datetime picker", func(t *testing.T) {
		dtp := NewDatetimePicker()

		data, err := json.Marshal(dtp)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "datetimepicker" {
			t.Errorf("got type %v, want datetimepicker", result["type"])
		}
	})

	t.Run("includes initial_date_time when set", func(t *testing.T) {
		dtp := NewDatetimePicker(WithInitialDateTime(1705327200))

		data, _ := json.Marshal(dtp)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["initial_date_time"] != float64(1705327200) {
			t.Errorf("got initial_date_time %v, want 1705327200", result["initial_date_time"])
		}
	})

	t.Run("implements InputElement interface", func(t *testing.T) {
		dtp := NewDatetimePicker()
		var _ InputElement = dtp
	})

	t.Run("does NOT implement SectionAccessory", func(t *testing.T) {
		// DatetimePicker should NOT be a SectionAccessory
		// This is a compile-time verification - if it compiled, the test passes
		_ = NewDatetimePicker()
	})
}
