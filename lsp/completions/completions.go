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
}
