package ast

import (
	"fmt"

	"github.com/textwire/textwire/v3/pkg/token"
)

type ContinueifDir struct {
	BaseNode
	Cond Expression
}

func NewContinueIfDir(tok token.Token) *ContinueifDir {
	return &ContinueifDir{
		BaseNode: NewBaseNode(tok),
	}
}

func (*ContinueifDir) chunkNode() {}

func (cd *ContinueifDir) String() string {
	return fmt.Sprintf("%s(%s)", cd.Token.Lit, cd.Cond)
}
