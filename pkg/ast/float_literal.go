package ast

import (
	"github.com/textwire/textwire/v3/pkg/token"
	"github.com/textwire/textwire/v3/pkg/utils"
)

type FloatLiteral struct {
	BaseNode
	Val float64
}

func NewFloatLiteral(tok token.Token, val float64) *FloatLiteral {
	return &FloatLiteral{
		BaseNode: NewBaseNode(tok),
		Val:      val,
	}
}

func (fl *FloatLiteral) expressionNode() {}

func (fl *FloatLiteral) String() string {
	return utils.FloatToStr(fl.Val)
}
