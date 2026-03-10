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

func (_ *IllegalNode) statementNode()  {}
func (_ *IllegalNode) expressionNode() {}
func (_ *IllegalNode) chunkNode()      {}

func (_ *IllegalNode) String() string {
	return ""
}
