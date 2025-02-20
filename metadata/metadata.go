package metadata

import (
	_ "embed"

	"github.com/textwire/textwire/v2/token"
)

//go:embed meta/if.md
var ifMeta string

//go:embed meta/ifelse.md
var elseIf string

var TokenMeta = map[token.TokenType]string{
	token.IF:      ifMeta,
	token.ELSE_IF: elseIf,
}

func GetTokenDoc(tok token.TokenType) string {
	meta, ok := TokenMeta[tok]
	if !ok {
		return ""
	}

	return meta
}
