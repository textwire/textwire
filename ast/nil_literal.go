package ast

import "github.com/textwire/textwire/v3/token"

type NilLiteral struct {
	BaseNode
}

func NewNilLiteral(tok token.Token) *NilLiteral {
	return &NilLiteral{
		BaseNode{
			Token: tok,
			Pos:   tok.Pos,
		},
	}
}

func (nl *NilLiteral) expressionNode() {}

func (nl *NilLiteral) String() string {
	return nl.Token.Literal
}
