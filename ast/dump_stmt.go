package ast

import (
	"bytes"

	"github.com/textwire/textwire/v2/token"
)

type DumpStmt struct {
	Token     token.Token // The '@dump' token
	Arguments []Expression
	Pos       Position
}

func (ds *DumpStmt) statementNode() {
}

func (ds *DumpStmt) TokenLiteral() string {
	return ds.Token.Literal
}

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

func (ss *DumpStmt) Line() uint {
	return ss.Token.Line
}

func (ss *DumpStmt) Position() Position {
	return ss.Pos
}
