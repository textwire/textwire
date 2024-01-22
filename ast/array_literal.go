package ast

import (
	"bytes"

	"github.com/textwire/textwire/token"
)

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {
}

func (al *ArrayLiteral) TokenLiteral() string {
	return al.Token.Literal
}

func (al *ArrayLiteral) String() string {
	var result bytes.Buffer

	result.WriteString("[")

	for _, el := range al.Elements {
		result.WriteString(el.String() + ", ")
	}

	if result.Len() > 1 {
		result.Truncate(result.Len() - 2)
	}

	result.WriteString("]")

	return result.String()
}
