package ast

import (
	"bytes"
	"strings"

	"github.com/textwire/textwire/v2/token"
)

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
	var out bytes.Buffer

	for _, s := range bs.Statements {
		_, isHTML := s.(*HTMLStmt)
		str := s.String()

		if isHTML {
			out.WriteString(s.String())
		} else if strings.HasPrefix(str, "@") {
			out.WriteString(s.String())
		} else {
			out.WriteString("{{ " + s.String() + " }}")
		}
	}

	return out.String()
}

func (bs *BlockStmt) Stmts() []Statement {
	res := make([]Statement, 0)

	if bs.Statements == nil {
		return []Statement{}
	}

	for _, stmt := range bs.Statements {
		if stmt == nil {
			continue
		}

		if s, ok := stmt.(NodeWithStatements); ok {
			res = append(res, s.(Statement))
			res = append(res, s.Stmts()...)
		}
	}

	return res
}
