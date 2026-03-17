package ast

import (
	"strings"

	"github.com/textwire/textwire/v4/pkg/token"
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

func (*DumpDir) chunkNode() {}

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
