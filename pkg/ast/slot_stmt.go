package ast

import (
	"strings"

	"github.com/textwire/textwire/v3/pkg/token"
)

type SlotStmt struct {
	BaseNode
	IsLocal   bool           // If the slot from the external comp or local
	CompName  string         // Component name
	isDefault bool           // Whether the slot is named or default
	name      *StringLiteral // Empty when @slot is default
	block     *BlockStmt     // Optional block statement, can be nil
}

func NewSlotStmt(tok token.Token, name *StringLiteral, compName string, isLocal bool) *SlotStmt {
	return &SlotStmt{
		BaseNode: NewBaseNode(tok),
		name:     name,
		CompName: compName,
		IsLocal:  isLocal,
	}
}

func (ss *SlotStmt) statementNode() {}

func (ss *SlotStmt) Name() *StringLiteral {
	return ss.name
}

func (ss *SlotStmt) IsDefault() bool {
	return ss.isDefault
}

func (ss *SlotStmt) SetIsDefault(val bool) {
	ss.isDefault = val
}

func (ss *SlotStmt) SetBlock(b *BlockStmt) {
	ss.block = b
}

func (ss *SlotStmt) Block() *BlockStmt {
	return ss.block
}

func (ss *SlotStmt) String() string {
	var out strings.Builder
	out.Grow(6)

	if ss.name.Value == "" {
		out.WriteString("@slot")
	} else {
		out.WriteString("@slot(")
		out.WriteString(ss.name.String())
		out.WriteString(")")
	}

	if ss.block != nil {
		out.WriteString(ss.block.String())
		out.WriteString("@end")
	}

	return out.String()
}

func (ss *SlotStmt) Stmts() []Statement {
	if ss.block == nil {
		return []Statement{}
	}

	return ss.block.Stmts()
}
