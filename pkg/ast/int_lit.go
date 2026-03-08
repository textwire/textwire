package ast

import "github.com/textwire/textwire/v3/pkg/token"

type IntLit struct {
	BaseNode
	Val int64
}

func NewIntLit(tok token.Token, val int64) *IntLit {
	return &IntLit{
		BaseNode: NewBaseNode(tok),
		Val:      val,
	}
}

func (il *IntLit) expressionNode() {}

func (il *IntLit) String() string {
	return il.Token.Lit
}
