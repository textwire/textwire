package ast

import (
	"github.com/textwire/textwire/v3/pkg/token"
)

type IllegalNode struct {
	BaseNode
	chunkKind ChunkKind
}

func NewIllegalNode(tok token.Token, chunkKind ChunkKind) *IllegalNode {
	return &IllegalNode{
		BaseNode:  NewBaseNode(tok),
		chunkKind: chunkKind,
	}
}

func (_ *IllegalNode) statementNode()  {}
func (_ *IllegalNode) expressionNode() {}
func (_ *IllegalNode) chunkNode()      {}

func (in *IllegalNode) Kind() ChunkKind {
	return in.chunkKind
}

func (_ *IllegalNode) String() string {
	return ""
}
