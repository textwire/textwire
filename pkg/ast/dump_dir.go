package ast

import (
	"strings"

	"github.com/textwire/textwire/v3/pkg/token"
)

type DumpDir struct {
	BaseNode
	Args []Expression
}

func NewDumpDir(tok token.Token, args []Expression) *DumpDir {
	return &DumpDir{
		BaseNode: NewBaseNode(tok),
		Args:     args,
	}
}

func (_ *DumpDir) chunkNode() {}

func (_ *DumpDir) Kind() ChunkKind {
	return ChunkKindDirective
}

func (dd *DumpDir) String() string {
	var out strings.Builder
	out.Grow(len(dd.Args) * 3)

	out.WriteString("@dump(")

	for i := range dd.Args {
		out.WriteString(dd.Args[i].String())

		if i < len(dd.Args)-1 {
			out.WriteString(",")
		}
	}

	out.WriteString(")")

	return out.String()
}
