package ast

import (
	"github.com/textwire/textwire/v3/pkg/token"
)

type IllegalNode struct {
	BaseNode
}

func NewIllegalNode(tok token.Token) *IllegalNode {
	return &IllegalNode{
		BaseNode: NewBaseNode(tok),
	}
}

func (*IllegalNode) statementNode()  {}
func (*IllegalNode) expressionNode() {}
func (*IllegalNode) chunkNode()      {}

func (*IllegalNode) String() string {
	return ""
}
