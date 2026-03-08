package ast

import (
	"github.com/textwire/textwire/v3/pkg/token"
	"github.com/textwire/textwire/v3/pkg/utils"
)

type FloatLit struct {
	BaseNode
	Val float64
}

func NewFloatLit(tok token.Token, val float64) *FloatLit {
	return &FloatLit{
		BaseNode: NewBaseNode(tok),
		Val:      val,
	}
}

func (fl *FloatLit) expressionNode() {}

func (fl *FloatLit) String() string {
	return utils.FloatToStr(fl.Val)
}
