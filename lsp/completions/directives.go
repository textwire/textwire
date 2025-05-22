package completions

import (
	"github.com/textwire/textwire/v2/lsp"
	"github.com/textwire/textwire/v2/lsp/utils"
	"github.com/textwire/textwire/v2/token"
)

func GetDirectives(locale lsp.Locale) ([]Completion, error) {
	directives := token.GetDirectives()
	completions := make([]Completion, 0, len(directives))

	for dir, tok := range directives {
		meta, err := lsp.GetTokenMeta(tok, locale)

		if err != nil {
			if err.Error() == utils.ErrNoMetadataFound {
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
