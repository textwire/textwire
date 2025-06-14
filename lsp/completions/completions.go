package completions

type Completion struct {
	// Label is text that is shown in the dropdown suggestions list.
	Label string

	// Insert is text that should be inserted into the editor after
	// the confirmation. Can be an empty string, in which case the
	// Label will be used for insertion.
	Insert string

	// Documentation is a full description of the item.
	Documentation string

	// InsertTextFormat defines whether the insert text in a completion
	// item should be interpreted as plain text or a snippet.
	// 1 stands for text format
	// 2 stands for snippet where you can include `${3:foo}` and just `$1`
	// special symbols to tell where to place the cursor
	InsertTextFormat int `json:"insertTextFormat"`
}
