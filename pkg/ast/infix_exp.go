package ast

import (
	"fmt"

	"github.com/textwire/textwire/v3/pkg/token"
)

type InfixExp struct {
	BaseNode
	Op    string // +, -, *, /, etc.
	Left  Expression
	Right Expression
}

func NewInfixExp(tok token.Token, left Expression, op string) *InfixExp {
	return &InfixExp{
		BaseNode: NewBaseNode(tok),
		Left:     left,
		Op:       op,
	}
}

func (ie *InfixExp) expressionNode() {}

func (ie *InfixExp) String() string {
	return fmt.Sprintf("(%s %s %s)", ie.Left, ie.Op, ie.Right)
}
