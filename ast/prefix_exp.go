package ast

import (
	"fmt"

	"github.com/textwire/textwire/v2/token"
)

type PrefixExp struct {
	BaseNode
	Operator string
	Right    Expression
}

func NewPrefixExp(tok token.Token, op string) *PrefixExp {
	return &PrefixExp{
		BaseNode: NewBaseNode(tok),
		Operator: op,
	}
}

func (pe *PrefixExp) expressionNode() {}

func (pe *PrefixExp) String() string {
	return fmt.Sprintf("(%s%s)", pe.Operator, pe.Right)
}
