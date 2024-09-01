package ast

import (
	"bytes"

	"github.com/textwire/textwire/token"
)

type ComponentStmt struct {
	Token    token.Token // The '@component' token
	Name     *StringLiteral
	Argument *ObjectLiteral
	Block    *Program
	Slots    []*SlotStmt
}

func (cs *ComponentStmt) statementNode() {
}

func (cs *ComponentStmt) TokenLiteral() string {
	return cs.Token.Literal
}

func (cs *ComponentStmt) String() string {
	var out bytes.Buffer

	out.WriteString("@component(")
	out.WriteString(cs.Name.String())
	out.WriteString(", ")
	out.WriteString(cs.Argument.String())
	out.WriteString(")")

	for _, slot := range cs.Slots {
		out.WriteString("\n")
		out.WriteString(slot.String())
	}

	return out.String()
}

func (cs *ComponentStmt) Line() uint {
	return cs.Token.Line
}
