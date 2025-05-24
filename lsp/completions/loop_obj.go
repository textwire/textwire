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
				Label:  "index",
				Insert: "index",
				Documentation: fmt.Sprintf("%s\n%s",
					"(property) index: int",
					"The current iteration of the loop. Starts with 0"),
			},
			{
				Label:  "first",
				Insert: "first",
				Documentation: fmt.Sprintf("%s\n%s",
					"(property) first: bool",
					"Returns `true` if this is the first iteration of the loop."),
			},
			{
				Label:  "last",
				Insert: "last",
				Documentation: fmt.Sprintf("%s\n%s",
					"(property) last: bool",
					"Returns `true` if this is the last iteration of the loop."),
			},
			{
				Label:  "iter",
				Insert: "iter",
				Documentation: fmt.Sprintf("%s\n%s",
					"(property) iter: int",
					"The current iteration of the loop. Starts with 1"),
			},
		},
	}

	_, ok := completions[locale]

	if !ok {
		return []Completion{}, fmt.Errorf(utils.ErrInvalidLocale, locale)
	}

	return completions[locale], nil
}
