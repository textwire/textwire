package completions

type Completion struct {
	// Label is text that is shown in the dropdown suggestions list.
	Label string

	// InsertText is a string that should be inserted into a document when
	// selecting this completion. When omitted the label is used as the
	// insert text for this item.
	InsertText string

	// Documentation is a full description of the item.
	Documentation string

	// InsertTextFormat defines whether the insert text in a completion
	// item should be interpreted as plain text or a snippet.
	// 1 stands for text format
	// 2 stands for snippet where you can include `${3:foo}` and just `$1`
	// special symbols to tell where to place the cursor
	InsertTextFormat int `json:"insertTextFormat"`
}
