package ast

import (
	"fmt"
	"strings"

	"github.com/textwire/textwire/v3/pkg/token"
)

type ElseIfDir struct {
	BaseNode
	Condition Expression
	Block     *Block // @elseif()<Block>@end
}

func NewElseIfDir(tok token.Token) *ElseIfDir {
	return &ElseIfDir{
		BaseNode: NewBaseNode(tok),
	}
}

func (_ *ElseIfDir) chunkNode() {}

func (_ *ElseIfDir) Kind() ChunkKind {
	return ChunkKindDirective
}

func (ed *ElseIfDir) String() string {
	var out strings.Builder

	fmt.Fprintf(&out, "@elseif(%s)\n", ed.Condition)
	out.WriteString(ed.Block.String())

	return out.String()
}

func (ed *ElseIfDir) AllChunks() []Chunk {
	if ed.Block == nil {
		return []Chunk{}
	}

	return ed.Block.AllChunks()
}
