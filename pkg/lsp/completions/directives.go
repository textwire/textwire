package completions

import (
	"errors"

	"github.com/textwire/textwire/v3/pkg/lsp"
	"github.com/textwire/textwire/v3/pkg/lsp/utils"
	"github.com/textwire/textwire/v3/pkg/token"
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
