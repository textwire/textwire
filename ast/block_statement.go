package ast

import (
	"bytes"

	"github.com/textwire/textwire/token"
)

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {
}

func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		_, isHtml := s.(*HTMLStatement)

		if isHtml {
			out.WriteString(s.String())
		} else {
			out.WriteString("{{ " + s.String() + " }}")
		}
	}

	return out.String()
}
