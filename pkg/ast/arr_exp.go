package ast

import (
	"bytes"

	"github.com/textwire/textwire/v3/pkg/token"
)

type ArrExpr struct {
	BaseNode
	Elements []Expression
}

func NewArrExpr(tok token.Token) *ArrExpr {
	return &ArrExpr{
		BaseNode: NewBaseNode(tok),
	}
}

func (*ArrExpr) expressionNode() {}

func (ae *ArrExpr) String() string {
	var out bytes.Buffer
	out.Grow(len(ae.Elements) + (2 * len(ae.Elements)))

	out.WriteString("[")

	for _, el := range ae.Elements {
		out.WriteString(el.String() + ", ")
	}

	if out.Len() > 1 {
		out.Truncate(out.Len() - 2)
	}

	out.WriteString("]")

	return out.String()
}
