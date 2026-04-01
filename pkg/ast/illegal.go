package ast

import (
	"github.com/textwire/textwire/v4/pkg/token"
)

type Illegal struct {
	BaseNode
}

func NewIllegalNode(tok token.Token) *Illegal {
	return &Illegal{
		BaseNode: NewBaseNode(tok),
	}
}

func (*Illegal) statementNode()  {}
func (*Illegal) expressionNode() {}
func (*Illegal) chunkNode()      {}
func (*Illegal) segmentNode()    {}

func (*Illegal) String() string {
	return ""
}
