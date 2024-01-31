package ast

import (
	"fmt"

	"github.com/textwire/textwire/token"
)

type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Operator string      // The operator, e.g. +
	Left     Expression
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}

func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", ie.Left, ie.Operator, ie.Right)
}

func (ie *InfixExpression) LineNum() uint {
	return ie.Token.Line
}
