package ast

import "github.com/textwire/textwire/v3/pkg/token"

type IntExpr struct {
	BaseNode
	Val int64
}

func NewIntExpr(tok token.Token, val int64) *IntExpr {
	return &IntExpr{
		BaseNode: NewBaseNode(tok),
		Val:      val,
	}
}

func (*IntExpr) expressionNode() {}
func (*IntExpr) segmentNode()    {}

func (ie *IntExpr) String() string {
	return ie.Token.Lit
}
