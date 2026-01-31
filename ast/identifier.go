package ast

import "github.com/textwire/textwire/v3/token"

type Identifier struct {
	BaseNode
	Name string
}

func NewIdentifier(tok token.Token, name string) *Identifier {
	return &Identifier{
		BaseNode: NewBaseNode(tok),
		Name:     name,
	}
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) String() string {
	return i.Name
}
