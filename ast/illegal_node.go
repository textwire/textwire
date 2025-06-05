package ast

import (
	"github.com/textwire/textwire/v2/token"
)

type IllegalNode struct {
	BaseNode
}

func NewIllegalNode(tok token.Token) *IllegalNode {
	return &IllegalNode{
		BaseNode: NewBaseNode(tok),
	}
}

func (en *IllegalNode) statementNode()  {}
func (en *IllegalNode) expressionNode() {}

func (en *IllegalNode) String() string {
	return ""
}
