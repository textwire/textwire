package ast

import (
	"bytes"

	"github.com/textwire/textwire/v2/token"
)

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
	Pos      token.Position
}

func (al *ArrayLiteral) expressionNode() {}

func (al *ArrayLiteral) Tok() *token.Token {
	return &al.Token
}

func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("[")

	for _, el := range al.Elements {
		out.WriteString(el.String() + ", ")
	}

	if out.Len() > 1 {
		out.Truncate(out.Len() - 2)
	}

	out.WriteString("]")

	return out.String()
}

func (al *ArrayLiteral) Line() uint {
	return al.Token.ErrorLine()
}

func (al *ArrayLiteral) Position() token.Position {
	return al.Pos
}
