package completions

type Completion struct {
	// Label is text that is shown in the dropdown suggestions list.
	Label string

	// Insert is text that should be inserted into the editor after
	// the confirmation.
	Insert string

	// Documentation is a full description of the item.
	Documentation string
}
