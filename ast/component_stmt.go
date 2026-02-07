package ast

import (
	"strings"

	"github.com/textwire/textwire/v3/token"
)

type ComponentStmt struct {
	BaseNode
	Name     *StringLiteral // Relative path to the component 'components/book'
	Argument *ObjectLiteral
	CompProg *Program    // AST node of the component file Name
	Slots    []*SlotStmt // Each slot of the component's body
}

func NewComponentStmt(tok token.Token) *ComponentStmt {
	return &ComponentStmt{
		BaseNode: NewBaseNode(tok),
	}
}

func (cs *ComponentStmt) statementNode() {}

func (cs *ComponentStmt) ArgsString() string {
	var out strings.Builder
	out.Grow(10)

	out.WriteString(cs.Name.String())

	if cs.Argument != nil {
		out.WriteString(", ")
		out.WriteString(cs.Argument.String())
	}

	return out.String()
}

func (cs *ComponentStmt) String() string {
	var out strings.Builder
	out.Grow(len(cs.Slots) + 20)

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
	if cs.CompProg == nil {
		return []Statement{}
	}

	return cs.CompProg.Stmts()
}
