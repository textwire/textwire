package ast

import "github.com/textwire/textwire/v3/pkg/token"

type IdentExpr struct {
	BaseNode
	Name string
}

func NewIdentExpr(tok token.Token, name string) *IdentExpr {
	return &IdentExpr{
		BaseNode: NewBaseNode(tok),
		Name:     name,
	}
}

func (*IdentExpr) expressionNode() {}

func (ie *IdentExpr) String() string {
	return ie.Name
}
