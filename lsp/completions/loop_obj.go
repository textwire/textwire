package completions

import (
	"fmt"

	"github.com/textwire/textwire/v2/lsp"
	"github.com/textwire/textwire/v2/lsp/utils"
)

func GetLoopObjFields(locale lsp.Locale) ([]Completion, error) {
	completions := map[lsp.Locale][]Completion{
		"en": {
			{
				Label:      "index",
				InsertText: "index",
				Documentation: fmt.Sprintf("%s\n%s",
					"(property) index: int",
					"The current iteration of the loop. Starts with 0"),
			},
			{
				Label:      "first",
				InsertText: "first",
				Documentation: fmt.Sprintf("%s\n%s",
					"(property) first: bool",
					"Returns `true` if this is the first iteration of the loop."),
			},
			{
				Label:      "last",
				InsertText: "last",
				Documentation: fmt.Sprintf("%s\n%s",
					"(property) last: bool",
					"Returns `true` if this is the last iteration of the loop."),
			},
			{
				Label:      "iter",
				InsertText: "iter",
				Documentation: fmt.Sprintf("%s\n%s",
					"(property) iter: int",
					"The current iteration of the loop. Starts with 1"),
			},
		},
	}

	_, ok := completions[locale]

	if !ok {
		return []Completion{}, utils.ErrInvalidLocale(string(locale))
	}

	return completions[locale], nil
}
