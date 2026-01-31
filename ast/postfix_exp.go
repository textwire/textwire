package ast

import (
	"fmt"

	"github.com/textwire/textwire/v3/token"
)

type PostfixExp struct {
	BaseNode
	Operator string
	Left     Expression
}

func NewPostfixExp(tok token.Token, left Expression, op string) *PostfixExp {
	return &PostfixExp{
		BaseNode: NewBaseNode(tok),
		Left:     left,
		Operator: op, // "++" or "--"
	}
}

func (pe *PostfixExp) expressionNode() {}

func (pe *PostfixExp) String() string {
	return fmt.Sprintf("(%s%s)", pe.Left, pe.Operator)
}
