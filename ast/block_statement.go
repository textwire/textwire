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
	var result bytes.Buffer

	for _, s := range bs.Statements {
		_, isHtml := s.(*HTMLStatement)

		if isHtml {
			result.WriteString(s.String())
		} else {
			result.WriteString("{{ " + s.String() + " }}")
		}
	}

	return result.String()
}
