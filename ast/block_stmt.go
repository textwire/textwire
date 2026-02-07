package ast

import (
	"strings"

	"github.com/textwire/textwire/v3/token"
)

// BlockStmt holds the body of statements like @each, @if, etc.
type BlockStmt struct {
	BaseNode
	Statements []Statement
}

func NewBlockStmt(tok token.Token) *BlockStmt {
	return &BlockStmt{
		BaseNode: NewBaseNode(tok),
	}
}

func (bs *BlockStmt) statementNode() {}

func (bs *BlockStmt) String() string {
	var out strings.Builder

	for _, s := range bs.Statements {
		str := s.String()

		if s.Tok().Type == token.HTML {
			out.WriteString(str)
		} else if strings.HasPrefix(str, "@") {
			out.WriteString(str)
		} else {
			out.WriteString("{{ " + str + " }}")
		}
	}

	return out.String()
}

func (bs *BlockStmt) Stmts() []Statement {
	if bs.Statements == nil {
		return []Statement{}
	}

	stmts := make([]Statement, 0, len(bs.Statements))

	for _, stmt := range bs.Statements {
		if stmt == nil {
			continue
		}

		if s, ok := stmt.(NodeWithStatements); ok {
			stmts = append(stmts, s.(Statement))
			stmts = append(stmts, s.Stmts()...)
		}
	}

	return stmts
}
