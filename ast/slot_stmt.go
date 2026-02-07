package ast

import (
	"strings"

	"github.com/textwire/textwire/v3/token"
)

const DefaultSlotName = "__DEFAULT_SLOT__"

type SlotStmt struct {
	BaseNode
	Name     *StringLiteral
	Block    *BlockStmt // Optional block statement, can be nil
	IsLocal  bool       // If the slot from the external comp or local
	CompName string     // Component name
}

func NewSlotStmt(tok token.Token, name *StringLiteral, compName string, isLocal bool) *SlotStmt {
	return &SlotStmt{
		BaseNode: NewBaseNode(tok),
		Name:     name,
		CompName: compName,
		IsLocal:  isLocal,
	}
}

func (ss *SlotStmt) statementNode() {}

func (ss *SlotStmt) String() string {
	var out strings.Builder
	out.Grow(6)

	if ss.Name.Value == "" {
		out.WriteString("@slot")
	} else {
		out.WriteString("@slot(")
		out.WriteString(ss.Name.String())
		out.WriteString(")")
	}

	if ss.Block != nil {
		out.WriteString("\n")
		out.WriteString(ss.Block.String())
		out.WriteString("\n@end")
	}

	return out.String()
}

func (ss *SlotStmt) Stmts() []Statement {
	if ss.Block == nil {
		return []Statement{}
	}

	return ss.Block.Stmts()
}
