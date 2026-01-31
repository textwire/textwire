package ast

import (
	"bytes"

	"github.com/textwire/textwire/v3/token"
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
	var out bytes.Buffer

	out.WriteString("@dump(")

	for i, arg := range ds.Arguments {
		out.WriteString(arg.String())

		if i < len(ds.Arguments)-1 {
			out.WriteString(",")
		}
	}

	out.WriteString(")")

	return out.String()
}
