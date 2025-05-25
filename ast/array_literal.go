package ast

import (
	"bytes"

	"github.com/textwire/textwire/v2/token"
)

type ArrayLiteral struct {
	BaseNode
	Elements []Expression
}

func NewArrayLiteral(tok token.Token) *ArrayLiteral {
	return &ArrayLiteral{
		BaseNode: NewBaseNode(tok),
	}
}

func (al *ArrayLiteral) expressionNode() {}

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
