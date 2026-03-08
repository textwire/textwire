package ast

import (
	"bytes"

	"github.com/textwire/textwire/v3/pkg/token"
)

type ArrLit struct {
	BaseNode
	Elements []Expression
}

func NewArrLit(tok token.Token) *ArrLit {
	return &ArrLit{
		BaseNode: NewBaseNode(tok),
	}
}

func (al *ArrLit) expressionNode() {}

func (al *ArrLit) String() string {
	var out bytes.Buffer
	out.Grow(len(al.Elements) + (2 * len(al.Elements)))

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
