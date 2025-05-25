package ast

import "github.com/textwire/textwire/v2/token"

type Identifier struct {
	BaseNode
	Value string
}

func NewIdentifier(tok token.Token, val string) *Identifier {
	return &Identifier{
		BaseNode: NewBaseNode(tok),
		Value:    val,
	}
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) String() string {
	return i.Value
}
