package ast

import "github.com/textwire/textwire/v2/token"

type StringLiteral struct {
	Token token.Token // The content of the string
	Value string
	Pos   token.Position
}

func NewStringLiteral(tok token.Token, val string) *StringLiteral {
	return &StringLiteral{
		Token: tok,
		Pos:   tok.Pos,
		Value: val,
	}
}

func (sl *StringLiteral) expressionNode() {}

func (sl *StringLiteral) Tok() *token.Token {
	return &sl.Token
}

func (sl *StringLiteral) String() string {
	return `"` + sl.Token.Literal + `"`
}

func (sl *StringLiteral) Line() uint {
	return sl.Token.ErrorLine()
}

func (sl *StringLiteral) Position() token.Position {
	return sl.Pos
}
