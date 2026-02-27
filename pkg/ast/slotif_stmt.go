package ast

import (
	"strings"

	"github.com/textwire/textwire/v3/pkg/token"
)

// SlotifStmt cannot be external and live in component file.
type SlotifStmt struct {
	BaseNode  // @slotif(bool, 'name'?)
	Condition Expression
	CompName  string         // Component name
	isDefault bool           // Whether the slot is named or default
	name      *StringLiteral // Empty when @slot is default
	block     *BlockStmt     // Required, cannot be nil
}

func NewSlotifStmt(tok token.Token, name *StringLiteral, compName string) *SlotifStmt {
	return &SlotifStmt{
		BaseNode: NewBaseNode(tok),
		name:     name,
		CompName: compName,
	}
}

func (sis *SlotifStmt) statementNode() {}

func (sis *SlotifStmt) Name() *StringLiteral {
	return sis.name
}

func (sis *SlotifStmt) IsDefault() bool {
	return sis.isDefault
}

func (sis *SlotifStmt) SetIsDefault(val bool) {
	sis.isDefault = val
}

func (sis *SlotifStmt) SetBlock(b *BlockStmt) {
	sis.block = b
}

func (sis *SlotifStmt) Block() *BlockStmt {
	return sis.block
}

func (sis *SlotifStmt) String() string {
	var out strings.Builder
	out.Grow(6)

	out.WriteString(sis.Token.Literal)
	out.WriteString("(")
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

func (sis *SlotifStmt) Stmts() []Statement {
	if sis.block == nil {
		return []Statement{}
	}

	return sis.block.Stmts()
}
