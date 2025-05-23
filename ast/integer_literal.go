package ast

import "github.com/textwire/textwire/v2/token"

type IntegerLiteral struct {
	Token token.Token
	Value int64
	Pos   token.Position
}

func NewIntegerLiteral(tok token.Token, val int64) *IntegerLiteral {
	return &IntegerLiteral{
		Token: tok,
		Pos:   tok.Pos,
		Value: val,
	}
}

func (il *IntegerLiteral) expressionNode() {}

func (il *IntegerLiteral) Tok() *token.Token {
	return &il.Token
}

func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

func (il *IntegerLiteral) Line() uint {
	return il.Token.ErrorLine()
}

func (il *IntegerLiteral) Position() token.Position {
	return il.Pos
}
