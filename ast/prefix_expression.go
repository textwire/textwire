package ast

import (
	"fmt"

	"github.com/textwire/textwire/token"
)

type PrefixExpression struct {
	Token    token.Token // The '!' or '-' token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {
}

func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", pe.Operator, pe.Right)
}

func (pe *PrefixExpression) Line() uint {
	return pe.Token.Line
}
