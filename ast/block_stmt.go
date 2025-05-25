package ast

import (
	"bytes"

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

		if isHTML {
			out.WriteString(s.String())
		} else {
			out.WriteString("{{ " + s.String() + " }}")
		}
	}

	return out.String()
}
