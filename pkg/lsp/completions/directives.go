package completions

import (
	"errors"

	"codeberg.org/textwire/textwire/v5/pkg/lsp"
	"codeberg.org/textwire/textwire/v5/pkg/lsp/utils"
	"codeberg.org/textwire/textwire/v5/pkg/token"
)

func GetDirectives(locale lsp.Locale) ([]Completion, error) {
	directives := token.GetDirectives()
	completions := make([]Completion, 0, len(directives))

	for dir, tok := range directives {
		meta, err := lsp.GetTokenMeta(tok, locale)
		if err != nil {
			if errors.Is(err, utils.ErrNoMetadataFound(tok)) {
				continue
			}
			return nil, err
		}

		insert, err := lsp.GetTokenInsert(tok)
		if err != nil {
			if errors.Is(err, utils.ErrNoMetadataFound(tok)) {
				continue
			}
			return nil, err
		}

		completions = append(completions, Completion{
			Label:            dir,
			InsertText:       insert[1:],
			InsertTextFormat: 2, // 2 = Snippet
			Documentation:    meta,
		})
	}

	return completions, nil
}
