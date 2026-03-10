package ast

import (
	"fmt"

	"github.com/textwire/textwire/v3/pkg/token"
)

type BreakifDir struct {
	BaseNode
	Cond Expression
}

func NewBreakIfDir(tok token.Token) *BreakifDir {
	return &BreakifDir{
		BaseNode: NewBaseNode(tok),
	}
}

func (_ *BreakifDir) chunkNode() {}

func (bd *BreakifDir) String() string {
	return fmt.Sprintf("%s(%s)", bd.Token.Lit, bd.Cond)
}
