package ast

import (
	"fmt"
	"strings"

	"github.com/textwire/textwire/v4/pkg/token"
)

type ElseIfDir struct {
	BaseNode
	Cond  Expression
	Block *Block // @elseif()<Block>@end
}

func NewElseIfDir(tok token.Token) *ElseIfDir {
	return &ElseIfDir{
		BaseNode: NewBaseNode(tok),
	}
}

func (*ElseIfDir) chunkNode() {}

func (ed *ElseIfDir) String() string {
	var out strings.Builder

	fmt.Fprintf(&out, "@elseif(%s)\n", ed.Cond)
	out.WriteString(ed.Block.String())

	return out.String()
}

func (ed *ElseIfDir) AllChunks() []Chunk {
	if ed.Block == nil {
		return []Chunk{}
	}

	return ed.Block.AllChunks()
}
