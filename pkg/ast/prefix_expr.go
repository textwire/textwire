package ast

import (
	"fmt"

	"github.com/textwire/textwire/v4/pkg/token"
)

type PrefixExpr struct {
	BaseNode
	Op    string // ! or -
	Right Expression
}

func NewPrefixExpr(tok token.Token, op string) *PrefixExpr {
	return &PrefixExpr{
		BaseNode: NewBaseNode(tok),
		Op:       op,
	}
}

func (*PrefixExpr) expressionNode() {}
func (*PrefixExpr) segmentNode()    {}

func (pe *PrefixExpr) String() string {
	return fmt.Sprintf("(%s%s)", pe.Op, pe.Right)
}
