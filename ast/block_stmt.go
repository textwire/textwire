package ast

import (
	"bytes"

	"github.com/textwire/textwire/token"
)

type BlockStmt struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStmt) statementNode() {
}

func (bs *BlockStmt) TokenLiteral() string {
	return bs.Token.Literal
}

func (bs *BlockStmt) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		_, isHTML := s.(*HTMLStatement)

		if isHTML {
			out.WriteString(s.String())
		} else {
			out.WriteString("{{ " + s.String() + " }}")
		}
	}

	return out.String()
}

func (bs *BlockStmt) Line() uint {
	return bs.Token.Line
}
