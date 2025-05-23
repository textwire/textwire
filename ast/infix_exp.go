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

func NewInfixExp(tok token.Token, left Expression, op string) *InfixExp {
	return &InfixExp{
		Token:    tok, // operator
		Pos:      tok.Pos,
		Left:     left,
		Operator: op,
	}
}

func (ie *InfixExp) expressionNode() {}

func (ie *InfixExp) Tok() *token.Token {
	return &ie.Token
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
