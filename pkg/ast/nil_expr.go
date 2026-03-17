package ast

import "github.com/textwire/textwire/v4/pkg/token"

type NilExpr struct {
	BaseNode
}

func NewNilExpr(tok token.Token) *NilExpr {
	return &NilExpr{NewBaseNode(tok)}
}

func (*NilExpr) expressionNode() {}
func (*NilExpr) segmentNode()    {}

func (ne *NilExpr) String() string {
	return ne.Token.Lit
}
