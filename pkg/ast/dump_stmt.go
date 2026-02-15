package ast

import (
	"strings"

	"github.com/textwire/textwire/v3/pkg/token"
)

type DumpStmt struct {
	BaseNode
	Arguments []Expression
}

func NewDumpStmt(tok token.Token, args []Expression) *DumpStmt {
	return &DumpStmt{
		BaseNode:  NewBaseNode(tok),
		Arguments: args,
	}
}

func (ds *DumpStmt) statementNode() {}

func (ds *DumpStmt) String() string {
	var out strings.Builder
	out.Grow(len(ds.Arguments) * 3)

	out.WriteString("@dump(")

	for i := range ds.Arguments {
		out.WriteString(ds.Arguments[i].String())

		if i < len(ds.Arguments)-1 {
			out.WriteString(",")
		}
	}

	out.WriteString(")")

	return out.String()
}
