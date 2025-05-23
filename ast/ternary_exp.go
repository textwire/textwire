package ast

import (
	"fmt"

	"github.com/textwire/textwire/v2/token"
)

type TernaryExp struct {
	Token       token.Token // The '?' token
	Condition   Expression
	Consequence Expression
	Alternative Expression
	Pos         token.Position
}

func NewTernaryExp(tok token.Token, cond Expression) *TernaryExp {
	return &TernaryExp{
		Token:     tok, // boolean condition
		Pos:       tok.Pos,
		Condition: cond,
	}
}

func (te *TernaryExp) expressionNode() {}

func (te *TernaryExp) Tok() *token.Token {
	return &te.Token
}

func (te *TernaryExp) String() string {
	return fmt.Sprintf("(%s ? %s : %s)", te.Condition, te.Condition, te.Alternative)
}

func (te *TernaryExp) Line() uint {
	return te.Token.ErrorLine()
}

func (te *TernaryExp) Position() token.Position {
	return te.Pos
}
