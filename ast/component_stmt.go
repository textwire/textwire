package ast

import (
	"bytes"

	"github.com/textwire/textwire/v2/token"
)

type ComponentStmt struct {
	BaseNode
	Name     *StringLiteral
	Argument *ObjectLiteral
	Block    *Program
	Slots    []*SlotStmt
}

func NewComponentStmt(tok token.Token) *ComponentStmt {
	return &ComponentStmt{
		BaseNode: NewBaseNode(tok),
	}
}

func (cs *ComponentStmt) statementNode() {}

func (cs *ComponentStmt) ArgsString() string {
	var out bytes.Buffer

	out.WriteString(cs.Name.String())

	if cs.Argument != nil {
		out.WriteString(", ")
		out.WriteString(cs.Argument.String())
	}

	return out.String()
}

func (cs *ComponentStmt) String() string {
	var out bytes.Buffer

	out.WriteString("@component(")
	out.WriteString(cs.ArgsString())
	out.WriteString(")")

	for _, slot := range cs.Slots {
		out.WriteString("\n")
		out.WriteString(slot.String())
	}

	if len(cs.Slots) > 0 {
		out.WriteString("\n@end\n")
	}

	return out.String()
}

func (cs *ComponentStmt) Stmts() []Statement {
	if cs.Block == nil {
		return []Statement{}
	}

	return cs.Block.Stmts()
}
