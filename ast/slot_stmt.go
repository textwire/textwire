package ast

import (
	"bytes"

	"github.com/textwire/textwire/v3/token"
)

type SlotStmt struct {
	BaseNode
	Name *StringLiteral // when empty string literal, it means default slot
	Body *BlockStmt     // optional block statement, can be nil
}

func NewSlotStmt(tok token.Token, name *StringLiteral) *SlotStmt {
	return &SlotStmt{
		BaseNode: NewBaseNode(tok),
		Name:     name,
	}
}

func (ss *SlotStmt) statementNode() {}

func (ss *SlotStmt) String() string {
	var out bytes.Buffer

	if ss.Name.Value == "" {
		out.WriteString("@slot")
	} else {
		out.WriteString("@slot(")
		out.WriteString(ss.Name.String())
		out.WriteString(")")
	}

	if ss.Body != nil {
		out.WriteString("\n")
		out.WriteString(ss.Body.String())
		out.WriteString("\n@end")
	}

	return out.String()
}

func (ss *SlotStmt) Stmts() []Statement {
	if ss.Body == nil {
		return []Statement{}
	}

	return ss.Body.Stmts()
}
