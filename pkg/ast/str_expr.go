package ast

import (
	"fmt"

	"github.com/textwire/textwire/v4/pkg/token"
)

type StrExpr struct {
	BaseNode
	Val string
}

func NewStrExpr(tok token.Token, val string) *StrExpr {
	return &StrExpr{
		BaseNode: NewBaseNode(tok),
		Val:      val,
	}
}

func (*StrExpr) expressionNode() {}
func (*StrExpr) segmentNode()    {}

func (se *StrExpr) String() string {
	return fmt.Sprintf(`"%s"`, se.Val)
}
