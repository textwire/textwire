package ast

import (
	"fmt"

	"github.com/textwire/textwire/v2/token"
)

type PostfixExp struct {
	Token    token.Token // The '++' or '--' token
	Operator string
	Left     Expression
	Pos      token.Position
}

func NewPostfixExp(tok token.Token, left Expression, op string) *PostfixExp {
	return &PostfixExp{
		Token:    tok, // identifier
		Pos:      tok.Pos,
		Left:     left,
		Operator: op, // "++" or "--"
	}
}

func (pe *PostfixExp) expressionNode() {}

func (pe *PostfixExp) Tok() *token.Token {
	return &pe.Token
}

func (pe *PostfixExp) String() string {
	return fmt.Sprintf("(%s%s)", pe.Left, pe.Operator)
}

func (pe *PostfixExp) Line() uint {
	return pe.Token.ErrorLine()
}

func (pe *PostfixExp) Position() token.Position {
	return pe.Pos
}
