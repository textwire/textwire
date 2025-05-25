package ast

import "github.com/textwire/textwire/v2/token"

type IntegerLiteral struct {
	BaseNode
	Value int64
}

func NewIntegerLiteral(tok token.Token, val int64) *IntegerLiteral {
	return &IntegerLiteral{
		BaseNode: NewBaseNode(tok),
		Value:    val,
	}
}

func (il *IntegerLiteral) expressionNode() {}

func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}
