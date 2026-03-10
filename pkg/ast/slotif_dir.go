package ast

import (
	"strings"

	"github.com/textwire/textwire/v3/pkg/token"
)

// SlotifDir cannot be external and live in component file.
type SlotifDir struct {
	BaseNode  // @slotif(bool, 'name'?)
	Cond      Expression
	CompName  string   // Component name
	isDefault bool     // Whether the slot is named or default
	name      *StrExpr // Empty when @slot is default
	block     *Block   // Required, cannot be nil
}

func NewSlotifDir(tok token.Token, name *StrExpr, compName string) *SlotifDir {
	return &SlotifDir{
		BaseNode: NewBaseNode(tok),
		CompName: compName,
		name:     name,
	}
}

func (_ *SlotifDir) chunkNode() {}

func (sd *SlotifDir) Name() *StrExpr {
	return sd.name
}

func (sd *SlotifDir) IsDefault() bool {
	return sd.isDefault
}

func (sd *SlotifDir) SetIsDefault(val bool) {
	sd.isDefault = val
}

func (sd *SlotifDir) SetBlock(b *Block) {
	sd.block = b
}

func (sd *SlotifDir) Block() *Block {
	return sd.block
}

func (sd *SlotifDir) String() string {
	var out strings.Builder
	out.Grow(6)

	out.WriteString(sd.Token.Lit)
	out.WriteString("(")
	out.WriteString(sd.Cond.String())

	if sd.name.Val != "" {
		out.WriteString(", ")
		out.WriteString(sd.name.String())
	}

	out.WriteString(")")
	out.WriteString(sd.block.String())
	out.WriteString("@end")

	return out.String()
}

func (sd *SlotifDir) AllChunks() []Chunk {
	if sd.block == nil {
		return []Chunk{}
	}
	return sd.block.AllChunks()
}
