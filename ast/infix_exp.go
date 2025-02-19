package ast

import (
	"fmt"

	"github.com/textwire/textwire/v2/token"
)

type InfixExp struct {
	Token    token.Token // The operator token, e.g. +
	Operator string      // The operator, e.g. +
	Left     Expression
	Right    Expression
	Pos      token.Position
}

func (ie *InfixExp) expressionNode() {}

func (ie *InfixExp) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *InfixExp) String() string {
	return fmt.Sprintf("(%s %s %s)", ie.Left, ie.Operator, ie.Right)
}

func (ie *InfixExp) Line() uint {
	return ie.Token.ErrorLine()
}

func (ie *InfixExp) Position() token.Position {
	return ie.Pos
}
