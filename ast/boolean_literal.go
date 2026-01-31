package ast

import (
	"github.com/textwire/textwire/v3/token"
)

type BooleanLiteral struct {
	BaseNode
	Value bool
}

func NewBooleanLiteral(tok token.Token, val bool) *BooleanLiteral {
	return &BooleanLiteral{
		BaseNode: NewBaseNode(tok),
		Value:    val,
	}
}

func (bl *BooleanLiteral) expressionNode() {}

func (bl *BooleanLiteral) String() string {
	return bl.Token.Literal
}
