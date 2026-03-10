package ast

import (
	"github.com/textwire/textwire/v3/pkg/token"
)

type BreakDir struct {
	BaseNode
}

func NewBreakDir(tok token.Token) *BreakDir {
	return &BreakDir{
		BaseNode: NewBaseNode(tok),
	}
}

func (_ *BreakDir) chunkNode() {}

func (bd *BreakDir) String() string {
	return bd.Token.Lit
}
