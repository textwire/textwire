package ast

import "github.com/textwire/textwire/v3/pkg/token"

type Ident struct {
	BaseNode
	Name string
}

func NewIdent(tok token.Token, name string) *Ident {
	return &Ident{
		BaseNode: NewBaseNode(tok),
		Name:     name,
	}
}

func (i *Ident) expressionNode() {}

func (i *Ident) String() string {
	return i.Name
}
