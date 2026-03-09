package ast

import (
	"strings"

	"github.com/textwire/textwire/v3/pkg/token"
)

type DumpDir struct {
	BaseNode
	Arguments []Expression
}

func NewDumpDir(tok token.Token, args []Expression) *DumpDir {
	return &DumpDir{
		BaseNode:  NewBaseNode(tok),
		Arguments: args,
	}
}

func (_ *DumpDir) chunkNode() {}

func (_ *DumpDir) Kind() ChunkKind {
	return ChunkKindDirective
}

func (dd *DumpDir) String() string {
	var out strings.Builder
	out.Grow(len(dd.Arguments) * 3)

	out.WriteString("@dump(")

	for i := range dd.Arguments {
		out.WriteString(dd.Arguments[i].String())

		if i < len(dd.Arguments)-1 {
			out.WriteString(",")
		}
	}

	out.WriteString(")")

	return out.String()
}
