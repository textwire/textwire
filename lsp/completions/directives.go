package completions

import (
	"github.com/textwire/textwire/v2/lsp"
	"github.com/textwire/textwire/v2/token"
)

func GetDirectives(locale lsp.Locale) ([]Completion, error) {
	completions := make([]Completion, 10)
	directives := token.GetDirectives()

	for dir, tok := range directives {
		meta, err := lsp.GetTokenMeta(tok, locale)

		if err != nil {
			if err.Error() == lsp.ErrNoMetadataFound {
				continue
			}

			return nil, err
		}

		completions = append(completions, Completion{
			Label:         dir,
			Insert:        dir[1:],
			Documentation: meta,
		})
	}

	return completions, nil
}
