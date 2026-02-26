package ast

import (
	"strings"

	"github.com/textwire/textwire/v3/pkg/token"
)

// SlotIfStmt cannot be external and live in component file.
type SlotIfStmt struct {
	BaseNode  // @slotIf(bool, 'name'?)
	Condition Expression
	CompName  string         // Component name
	isDefault bool           // Whether the slot is named or default
	name      *StringLiteral // Empty when @slot is default
	block     *BlockStmt     // Required, cannot be nil
}

func NewSlotIfStmt(tok token.Token, name *StringLiteral, compName string) *SlotIfStmt {
	return &SlotIfStmt{
		BaseNode: NewBaseNode(tok),
		name:     name,
		CompName: compName,
	}
}

func (sis *SlotIfStmt) statementNode() {}

func (sis *SlotIfStmt) Name() *StringLiteral {
	return sis.name
}

func (sis *SlotIfStmt) IsDefault() bool {
	return sis.isDefault
}

func (sis *SlotIfStmt) SetIsDefault(val bool) {
	sis.isDefault = val
}

func (sis *SlotIfStmt) SetBlock(b *BlockStmt) {
	sis.block = b
}

func (sis *SlotIfStmt) Block() *BlockStmt {
	return sis.block
}

func (sis *SlotIfStmt) String() string {
	var out strings.Builder
	out.Grow(6)

	out.WriteString("@slotIf(")
	out.WriteString(sis.Condition.String())

	if sis.name.Value != "" {
		out.WriteString(", ")
		out.WriteString(sis.name.String())
	}

	out.WriteString(")")
	out.WriteString(sis.block.String())
	out.WriteString("@end")

	return out.String()
}

func (sis *SlotIfStmt) Stmts() []Statement {
	if sis.block == nil {
		return []Statement{}
	}

	return sis.block.Stmts()
}
