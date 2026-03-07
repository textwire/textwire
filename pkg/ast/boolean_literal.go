package ast

import (
	"github.com/textwire/textwire/v3/pkg/token"
)

type BooleanLiteral struct {
	BaseNode
	Val bool
}

func NewBooleanLiteral(tok token.Token, val bool) *BooleanLiteral {
	return &BooleanLiteral{
		BaseNode: NewBaseNode(tok),
		Val:      val,
	}
}

func (bl *BooleanLiteral) expressionNode() {}

func (bl *BooleanLiteral) String() string {
	return bl.Token.Literal
}
