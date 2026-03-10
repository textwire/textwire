package ast

import "github.com/textwire/textwire/v3/pkg/token"

type NilExpr struct {
	BaseNode
}

func NewNilExpr(tok token.Token) *NilExpr {
	return &NilExpr{
		BaseNode{
			Token: tok,
			Pos:   tok.Pos,
		},
	}
}

func (*NilExpr) expressionNode() {}

func (ne *NilExpr) String() string {
	return ne.Token.Lit
}
