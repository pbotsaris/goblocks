package blocks

import "encoding/json"

// Table displays structured information in a table format.
// Available in messages only.
type Table struct {
	columns []TableColumn
	rows    []TableRow
	blockID string
}

// Marker interface implementation
func (Table) block() {}

// MarshalJSON implements json.Marshaler.
func (t Table) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type":    "table",
		"columns": t.columns,
		"rows":    t.rows,
	}
	if t.blockID != "" {
		m["block_id"] = t.blockID
	}
	return json.Marshal(m)
}

// TableOption configures a Table block.
type TableOption func(*Table)

// NewTable creates a new table block.
func NewTable(columns []TableColumn, rows []TableRow, opts ...TableOption) (Table, error) {
	if err := validateMinItems("columns", columns, 1); err != nil {
		return Table{}, err
	}
	if err := validateMinItems("rows", rows, 1); err != nil {
		return Table{}, err
	}

	t := Table{
		columns: columns,
		rows:    rows,
	}

	for _, opt := range opts {
		opt(&t)
	}

	return t, nil
}

// MustTable creates a Table or panics on error.
func MustTable(columns []TableColumn, rows []TableRow, opts ...TableOption) Table {
	t, err := NewTable(columns, rows, opts...)
	if err != nil {
		panic(err)
	}
	return t
}

// WithTableBlockID sets the block_id.
func WithTableBlockID(id string) TableOption {
	return func(t *Table) {
		t.blockID = id
	}
}

// TableColumn defines a column in a table.
type TableColumn struct {
	id    string
	name  string
	width int
}

// MarshalJSON implements json.Marshaler.
func (c TableColumn) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"id":   c.id,
		"name": c.name,
	}
	if c.width > 0 {
		m["width"] = c.width
	}
	return json.Marshal(m)
}

// TableColumnOption configures a TableColumn.
type TableColumnOption func(*TableColumn)

// NewTableColumn creates a new table column.
func NewTableColumn(id, name string, opts ...TableColumnOption) TableColumn {
	c := TableColumn{id: id, name: name}
	for _, opt := range opts {
		opt(&c)
	}
	return c
}

// WithTableColumnWidth sets the column width.
func WithTableColumnWidth(width int) TableColumnOption {
	return func(c *TableColumn) {
		c.width = width
	}
}

// TableRow represents a row in a table.
type TableRow struct {
	cells []TableCell
}

// MarshalJSON implements json.Marshaler.
func (r TableRow) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"cells": r.cells,
	})
}

// NewTableRow creates a new table row.
func NewTableRow(cells []TableCell) TableRow {
	return TableRow{cells: cells}
}

// TableCell represents a cell in a table row.
type TableCell struct {
	columnID string
	value    TableCellValue
}

// MarshalJSON implements json.Marshaler.
func (c TableCell) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"column_id": c.columnID,
		"value":     c.value,
	})
}

// NewTableCell creates a new table cell.
func NewTableCell(columnID string, value TableCellValue) TableCell {
	return TableCell{columnID: columnID, value: value}
}

// TableCellValue represents the value content of a table cell.
type TableCellValue interface {
	json.Marshaler
	tableCellValue()
}

// TableCellText is a text value for a table cell.
type TableCellText struct {
	text string
}

func (TableCellText) tableCellValue() {}

// MarshalJSON implements json.Marshaler.
func (t TableCellText) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type": "plain_text",
		"text": t.text,
	})
}

// NewTableCellText creates a new text value for a table cell.
func NewTableCellText(text string) TableCellText {
	return TableCellText{text: text}
}

// TableCellRichText is a rich text value for a table cell.
type TableCellRichText struct {
	elements []RichTextElement
}

func (TableCellRichText) tableCellValue() {}

// MarshalJSON implements json.Marshaler.
func (r TableCellRichText) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":     "rich_text",
		"elements": r.elements,
	})
}

// NewTableCellRichText creates a new rich text value for a table cell.
func NewTableCellRichText(elements []RichTextElement) TableCellRichText {
	return TableCellRichText{elements: elements}
}
