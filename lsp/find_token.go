package lsp

import (
	"github.com/textwire/textwire/v2/ast"
	"github.com/textwire/textwire/v2/token"
)

// FindToken recursively searches for a token at a given line and character
func FindToken(stmts []ast.Statement, line, char uint) (*token.Token, error) {
	for _, stmt := range stmts {
		pos := stmt.Position()

		// check if the current stmt is StatementsContainer interface
		if container, ok := stmt.(ast.StatementsContainer); ok {
			token, err := FindToken(container.Stmts(), line, char)
			if err != nil {
				return nil, err
			}

			if token != nil {
				return token, nil
			}
		}

		lineMatch := line >= pos.StartLine && line <= pos.EndLine
		colMatch := char >= pos.StartCol && char <= pos.EndCol
		if lineMatch && colMatch {
			return stmt.Tok(), nil
		}
	}

	return nil, nil
}
