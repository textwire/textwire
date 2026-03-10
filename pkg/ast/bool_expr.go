package ast

import (
	"github.com/textwire/textwire/v3/pkg/token"
)

type BoolExpr struct {
	BaseNode
	Val bool
}

func NewBoolExpr(tok token.Token, val bool) *BoolExpr {
	return &BoolExpr{
		BaseNode: NewBaseNode(tok),
		Val:      val,
	}
}

func (*BoolExpr) expressionNode() {}

func (be *BoolExpr) String() string {
	return be.Token.Lit
}
