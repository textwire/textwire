package ast

import (
	"fmt"

	"github.com/textwire/textwire/v4/pkg/token"
)

type InfixExpr struct {
	BaseNode
	Op    string // +, -, *, /, etc.
	Left  Expression
	Right Expression
}

func NewInfixExpr(tok token.Token, left Expression, op string) *InfixExpr {
	return &InfixExpr{
		BaseNode: NewBaseNode(tok),
		Left:     left,
		Op:       op,
	}
}

func (*InfixExpr) expressionNode() {}
func (*InfixExpr) segmentNode()    {}

func (ie *InfixExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", ie.Left, ie.Op, ie.Right)
}
