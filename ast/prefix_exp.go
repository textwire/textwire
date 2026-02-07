package ast

import (
	"fmt"

	"github.com/textwire/textwire/v3/token"
)

type PrefixExp struct {
	BaseNode
	Op    string
	Right Expression
}

func NewPrefixExp(tok token.Token, op string) *PrefixExp {
	return &PrefixExp{
		BaseNode: NewBaseNode(tok),
		Op:       op,
	}
}

func (pe *PrefixExp) expressionNode() {}

func (pe *PrefixExp) String() string {
	return fmt.Sprintf("(%s%s)", pe.Op, pe.Right)
}
