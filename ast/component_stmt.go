package ast

import (
	"bytes"

	"github.com/textwire/textwire/v2/token"
)

type ComponentStmt struct {
	Token    token.Token // The '@component' token
	Name     *StringLiteral
	Argument *ObjectLiteral
	Block    *Program
	Slots    []*SlotStmt
	Pos      token.Position
}

func (cs *ComponentStmt) statementNode() {}

func (cs *ComponentStmt) Stmts() []Statement {
	return cs.Block.Statements
}

func (cs *ComponentStmt) Tok() *token.Token {
	return &cs.Token
}

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

func (cs *ComponentStmt) Line() uint {
	return cs.Token.ErrorLine()
}

func (cs *ComponentStmt) Position() token.Position {
	return cs.Pos
}
