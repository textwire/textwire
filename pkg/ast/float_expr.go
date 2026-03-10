package ast

import (
	"github.com/textwire/textwire/v3/pkg/token"
	"github.com/textwire/textwire/v3/pkg/utils"
)

type FloatExpr struct {
	BaseNode
	Val float64
}

func NewFloatExpr(tok token.Token, val float64) *FloatExpr {
	return &FloatExpr{
		BaseNode: NewBaseNode(tok),
		Val:      val,
	}
}

func (*FloatExpr) expressionNode() {}
func (*FloatExpr) segmentNode()    {}

func (fe *FloatExpr) String() string {
	return utils.FloatToStr(fe.Val)
}
