package ast

import (
	"fmt"

	"github.com/textwire/textwire/v3/pkg/token"
)

type PostfixExp struct {
	BaseNode
	Op   string // ++ or --
	Left Expression
}

func NewPostfixExp(tok token.Token, left Expression, op string) *PostfixExp {
	return &PostfixExp{
		BaseNode: NewBaseNode(tok),
		Left:     left,
		Op:       op,
	}
}

func (pe *PostfixExp) expressionNode() {}

func (pe *PostfixExp) String() string {
	return fmt.Sprintf("(%s%s)", pe.Left, pe.Op)
}
