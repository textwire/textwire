package ast

import "github.com/textwire/textwire/v3/token"

type StringLiteral struct {
	BaseNode
	Value string
}

func NewStringLiteral(tok token.Token, val string) *StringLiteral {
	return &StringLiteral{
		BaseNode: NewBaseNode(tok),
		Value:    val,
	}
}

func (sl *StringLiteral) expressionNode() {}

func (sl *StringLiteral) String() string {
	return `"` + sl.Token.Literal + `"`
}
