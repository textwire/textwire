package ast

import (
	"github.com/textwire/textwire/v3/pkg/token"
)

type BoolLit struct {
	BaseNode
	Val bool
}

func NewBoolLit(tok token.Token, val bool) *BoolLit {
	return &BoolLit{
		BaseNode: NewBaseNode(tok),
		Val:      val,
	}
}

func (bl *BoolLit) expressionNode() {}

func (bl *BoolLit) String() string {
	return bl.Token.Lit
}
