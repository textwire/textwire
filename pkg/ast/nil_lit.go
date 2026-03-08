package ast

import "github.com/textwire/textwire/v3/pkg/token"

type NilLit struct {
	BaseNode
}

func NewNilLit(tok token.Token) *NilLit {
	return &NilLit{
		BaseNode{
			Token: tok,
			Pos:   tok.Pos,
		},
	}
}

func (nl *NilLit) expressionNode() {}

func (nl *NilLit) String() string {
	return nl.Token.Lit
}
