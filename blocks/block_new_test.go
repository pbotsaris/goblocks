package blocks

import (
	"encoding/json"
	"testing"
)

func TestFileBlock(t *testing.T) {
	t.Run("creates valid file block", func(t *testing.T) {
		file, err := NewFile("ABCD1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(file)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "file" {
			t.Errorf("got type %v, want 'file'", result["type"])
		}
		if result["external_id"] != "ABCD1" {
			t.Errorf("got external_id %v, want 'ABCD1'", result["external_id"])
		}
		if result["source"] != "remote" {
			t.Errorf("got source %v, want 'remote'", result["source"])
		}
	})

	t.Run("includes block_id when set", func(t *testing.T) {
		file, _ := NewFile("ABCD1", WithFileBlockID("file_1"))

		data, _ := json.Marshal(file)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["block_id"] != "file_1" {
			t.Errorf("got block_id %v, want 'file_1'", result["block_id"])
		}
	})

	t.Run("rejects empty external_id", func(t *testing.T) {
		_, err := NewFile("")
		if err == nil {
			t.Error("expected error for empty external_id")
		}
	})

	t.Run("implements Block interface", func(t *testing.T) {
		file, _ := NewFile("ABCD1")
		var _ Block = file
	})
}

func TestVideoBlock(t *testing.T) {
	t.Run("creates valid video block", func(t *testing.T) {
		video, err := NewVideo("Alt text", "Video Title",
			"https://example.com/thumb.png", "https://example.com/video.mp4")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(video)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "video" {
			t.Errorf("got type %v, want 'video'", result["type"])
		}
		if result["alt_text"] != "Alt text" {
			t.Errorf("got alt_text %v, want 'Alt text'", result["alt_text"])
		}
		if result["video_url"] != "https://example.com/video.mp4" {
			t.Errorf("got video_url %v, want 'https://example.com/video.mp4'", result["video_url"])
		}
	})

	t.Run("includes all options when set", func(t *testing.T) {
		video, _ := NewVideo("Alt", "Title",
			"https://example.com/thumb.png", "https://example.com/video.mp4",
			WithVideoAuthorName("John Doe"),
			WithVideoProviderName("YouTube"),
			WithVideoDescription("A great video"),
		)

		data, _ := json.Marshal(video)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["author_name"] != "John Doe" {
			t.Errorf("got author_name %v, want 'John Doe'", result["author_name"])
		}
		if result["provider_name"] != "YouTube" {
			t.Errorf("got provider_name %v, want 'YouTube'", result["provider_name"])
		}
	})

	t.Run("rejects empty alt_text", func(t *testing.T) {
		_, err := NewVideo("", "Title", "https://example.com/thumb.png", "https://example.com/video.mp4")
		if err == nil {
			t.Error("expected error for empty alt_text")
		}
	})

	t.Run("implements Block interface", func(t *testing.T) {
		video, _ := NewVideo("Alt", "Title", "https://example.com/thumb.png", "https://example.com/video.mp4")
		var _ Block = video
	})
}

func TestMarkdownBlock(t *testing.T) {
	t.Run("creates valid markdown block", func(t *testing.T) {
		md, err := NewMarkdownBlock("# Hello World")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(md)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "markdown" {
			t.Errorf("got type %v, want 'markdown'", result["type"])
		}
		if result["text"] != "# Hello World" {
			t.Errorf("got text %v, want '# Hello World'", result["text"])
		}
	})

	t.Run("includes block_id when set", func(t *testing.T) {
		md, _ := NewMarkdownBlock("# Hello", WithMarkdownBlockID("md_1"))

		data, _ := json.Marshal(md)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["block_id"] != "md_1" {
			t.Errorf("got block_id %v, want 'md_1'", result["block_id"])
		}
	})

	t.Run("rejects empty text", func(t *testing.T) {
		_, err := NewMarkdownBlock("")
		if err == nil {
			t.Error("expected error for empty text")
		}
	})

	t.Run("implements Block interface", func(t *testing.T) {
		md, _ := NewMarkdownBlock("# Hello")
		var _ Block = md
	})
}

func TestContextActionsBlock(t *testing.T) {
	t.Run("creates valid context actions block", func(t *testing.T) {
		fb := NewFeedbackButtons()
		ca, err := NewContextActions([]ContextActionsElement{fb})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(ca)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "context_actions" {
			t.Errorf("got type %v, want 'context_actions'", result["type"])
		}
		elements := result["elements"].([]any)
		if len(elements) != 1 {
			t.Errorf("got %d elements, want 1", len(elements))
		}
	})

	t.Run("includes block_id when set", func(t *testing.T) {
		fb := NewFeedbackButtons()
		ca, _ := NewContextActions([]ContextActionsElement{fb},
			WithContextActionsBlockID("ca_1"))

		data, _ := json.Marshal(ca)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["block_id"] != "ca_1" {
			t.Errorf("got block_id %v, want 'ca_1'", result["block_id"])
		}
	})

	t.Run("rejects empty elements", func(t *testing.T) {
		_, err := NewContextActions([]ContextActionsElement{})
		if err == nil {
			t.Error("expected error for empty elements")
		}
	})

	t.Run("implements Block interface", func(t *testing.T) {
		fb := NewFeedbackButtons()
		ca, _ := NewContextActions([]ContextActionsElement{fb})
		var _ Block = ca
	})
}

func TestRichTextBlock(t *testing.T) {
	t.Run("creates valid rich text block", func(t *testing.T) {
		section := NewRichTextSection([]RichTextSectionElement{
			NewRichTextText("Hello world", nil),
		})
		rt, err := NewRichText([]RichTextElement{section})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(rt)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "rich_text" {
			t.Errorf("got type %v, want 'rich_text'", result["type"])
		}
		elements := result["elements"].([]any)
		if len(elements) != 1 {
			t.Errorf("got %d elements, want 1", len(elements))
		}
	})

	t.Run("includes block_id when set", func(t *testing.T) {
		section := NewRichTextSection([]RichTextSectionElement{
			NewRichTextText("Hello", nil),
		})
		rt, _ := NewRichText([]RichTextElement{section},
			WithRichTextBlockID("rt_1"))

		data, _ := json.Marshal(rt)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["block_id"] != "rt_1" {
			t.Errorf("got block_id %v, want 'rt_1'", result["block_id"])
		}
	})

	t.Run("rejects empty elements", func(t *testing.T) {
		_, err := NewRichText([]RichTextElement{})
		if err == nil {
			t.Error("expected error for empty elements")
		}
	})

	t.Run("implements Block interface", func(t *testing.T) {
		section := NewRichTextSection([]RichTextSectionElement{
			NewRichTextText("Hello", nil),
		})
		rt, _ := NewRichText([]RichTextElement{section})
		var _ Block = rt
	})
}

func TestRichTextSection(t *testing.T) {
	t.Run("creates section with styled text", func(t *testing.T) {
		style := NewRichTextStyle(true, true, false, false)
		section := NewRichTextSection([]RichTextSectionElement{
			NewRichTextText("Bold and italic", style),
		})

		data, err := json.Marshal(section)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "rich_text_section" {
			t.Errorf("got type %v, want 'rich_text_section'", result["type"])
		}
	})

	t.Run("implements RichTextElement interface", func(t *testing.T) {
		section := NewRichTextSection([]RichTextSectionElement{
			NewRichTextText("Hello", nil),
		})
		var _ RichTextElement = section
	})
}

func TestRichTextList(t *testing.T) {
	t.Run("creates bullet list", func(t *testing.T) {
		section1 := NewRichTextSection([]RichTextSectionElement{
			NewRichTextText("Item 1", nil),
		})
		section2 := NewRichTextSection([]RichTextSectionElement{
			NewRichTextText("Item 2", nil),
		})
		list := NewRichTextList("bullet", []RichTextSection{section1, section2})

		data, err := json.Marshal(list)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "rich_text_list" {
			t.Errorf("got type %v, want 'rich_text_list'", result["type"])
		}
		if result["style"] != "bullet" {
			t.Errorf("got style %v, want 'bullet'", result["style"])
		}
	})

	t.Run("implements RichTextElement interface", func(t *testing.T) {
		section := NewRichTextSection([]RichTextSectionElement{
			NewRichTextText("Item", nil),
		})
		list := NewRichTextList("ordered", []RichTextSection{section})
		var _ RichTextElement = list
	})
}

func TestTableBlock(t *testing.T) {
	t.Run("creates valid table block", func(t *testing.T) {
		col1 := NewTableColumn("col1", "Name")
		col2 := NewTableColumn("col2", "Value")
		cell1 := NewTableCell("col1", NewTableCellText("Hello"))
		cell2 := NewTableCell("col2", NewTableCellText("World"))
		row := NewTableRow([]TableCell{cell1, cell2})

		table, err := NewTable([]TableColumn{col1, col2}, []TableRow{row})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(table)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "table" {
			t.Errorf("got type %v, want 'table'", result["type"])
		}
		columns := result["columns"].([]any)
		if len(columns) != 2 {
			t.Errorf("got %d columns, want 2", len(columns))
		}
	})

	t.Run("includes block_id when set", func(t *testing.T) {
		col := NewTableColumn("col1", "Name")
		cell := NewTableCell("col1", NewTableCellText("Hello"))
		row := NewTableRow([]TableCell{cell})
		table, _ := NewTable([]TableColumn{col}, []TableRow{row},
			WithTableBlockID("table_1"))

		data, _ := json.Marshal(table)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["block_id"] != "table_1" {
			t.Errorf("got block_id %v, want 'table_1'", result["block_id"])
		}
	})

	t.Run("implements Block interface", func(t *testing.T) {
		col := NewTableColumn("col1", "Name")
		cell := NewTableCell("col1", NewTableCellText("Hello"))
		row := NewTableRow([]TableCell{cell})
		table, _ := NewTable([]TableColumn{col}, []TableRow{row})
		var _ Block = table
	})
}

func TestPlanBlock(t *testing.T) {
	t.Run("creates valid plan block", func(t *testing.T) {
		item := MustPlanItem("Task 1", PlanItemStatusPending)
		section := MustPlanSection("Section 1", []PlanItem{item})
		plan, err := NewPlan("My Plan", []PlanSection{section})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(plan)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "plan" {
			t.Errorf("got type %v, want 'plan'", result["type"])
		}
	})

	t.Run("includes block_id when set", func(t *testing.T) {
		item := MustPlanItem("Task 1", PlanItemStatusComplete)
		section := MustPlanSection("Section 1", []PlanItem{item})
		plan, _ := NewPlan("My Plan", []PlanSection{section},
			WithPlanBlockID("plan_1"))

		data, _ := json.Marshal(plan)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["block_id"] != "plan_1" {
			t.Errorf("got block_id %v, want 'plan_1'", result["block_id"])
		}
	})

	t.Run("rejects empty title", func(t *testing.T) {
		item := MustPlanItem("Task 1", PlanItemStatusPending)
		section := MustPlanSection("Section 1", []PlanItem{item})
		_, err := NewPlan("", []PlanSection{section})
		if err == nil {
			t.Error("expected error for empty title")
		}
	})

	t.Run("implements Block interface", func(t *testing.T) {
		item := MustPlanItem("Task 1", PlanItemStatusPending)
		section := MustPlanSection("Section 1", []PlanItem{item})
		plan, _ := NewPlan("My Plan", []PlanSection{section})
		var _ Block = plan
	})
}

func TestTaskCardBlock(t *testing.T) {
	t.Run("creates valid task card block", func(t *testing.T) {
		tc, err := NewTaskCard("task_1", "My Task", TaskCardStatusOpen)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(tc)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "task_card" {
			t.Errorf("got type %v, want 'task_card'", result["type"])
		}
		if result["task_id"] != "task_1" {
			t.Errorf("got task_id %v, want 'task_1'", result["task_id"])
		}
		if result["status"] != "open" {
			t.Errorf("got status %v, want 'open'", result["status"])
		}
	})

	t.Run("includes all options when set", func(t *testing.T) {
		source := MustURLSource("https://example.com", WithURLSourceTitle("Example"))
		tc, _ := NewTaskCard("task_1", "My Task", TaskCardStatusInProgress,
			WithTaskCardDescription("A description"),
			WithTaskCardSources([]URLSource{source}),
			WithTaskCardBlockID("tc_1"),
		)

		data, _ := json.Marshal(tc)
		var result map[string]any
		mustUnmarshal(data, &result)

		if result["block_id"] != "tc_1" {
			t.Errorf("got block_id %v, want 'tc_1'", result["block_id"])
		}
		sources := result["sources"].([]any)
		if len(sources) != 1 {
			t.Errorf("got %d sources, want 1", len(sources))
		}
	})

	t.Run("rejects empty task_id", func(t *testing.T) {
		_, err := NewTaskCard("", "My Task", TaskCardStatusOpen)
		if err == nil {
			t.Error("expected error for empty task_id")
		}
	})

	t.Run("implements Block interface", func(t *testing.T) {
		tc, _ := NewTaskCard("task_1", "My Task", TaskCardStatusOpen)
		var _ Block = tc
	})
}

func TestURLSource(t *testing.T) {
	t.Run("creates valid URL source", func(t *testing.T) {
		source, err := NewURLSource("https://example.com")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := json.Marshal(source)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var result map[string]any
		mustUnmarshal(data, &result)

		if result["type"] != "url_source" {
			t.Errorf("got type %v, want 'url_source'", result["type"])
		}
		if result["url"] != "https://example.com" {
			t.Errorf("got url %v, want 'https://example.com'", result["url"])
		}
	})

	t.Run("includes title when set", func(t *testing.T) {
		source, _ := NewURLSource("https://example.com",
			WithURLSourceTitle("Example Site"))

		data, _ := json.Marshal(source)
		var result map[string]any
		mustUnmarshal(data, &result)

		title := result["title"].(map[string]any)
		if title["text"] != "Example Site" {
			t.Errorf("got title %v, want 'Example Site'", title["text"])
		}
	})

	t.Run("rejects empty URL", func(t *testing.T) {
		_, err := NewURLSource("")
		if err == nil {
			t.Error("expected error for empty URL")
		}
	})
}
