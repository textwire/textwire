package ast

import (
	"strings"

	"github.com/textwire/textwire/v3/pkg/token"
)

type SlotDir struct {
	BaseNode
	IsLocal   bool     // If the slot from the external comp or local
	CompName  string   // Component name
	isDefault bool     // Whether the slot is named or default
	name      *StrExpr // Empty when @slot is default
	block     *Block   // Optional block statement, can be nil
}

func NewSlotDir(tok token.Token, name *StrExpr, compName string, isLocal bool) *SlotDir {
	return &SlotDir{
		BaseNode: NewBaseNode(tok),
		name:     name,
		CompName: compName,
		IsLocal:  isLocal,
	}
}

func (_ *SlotDir) chunkNode() {}

func (sd *SlotDir) Name() *StrExpr {
	return sd.name
}

func (sd *SlotDir) IsDefault() bool {
	return sd.isDefault
}

func (sd *SlotDir) SetIsDefault(val bool) {
	sd.isDefault = val
}

func (sd *SlotDir) SetBlock(b *Block) {
	sd.block = b
}

func (sd *SlotDir) Block() *Block {
	return sd.block
}

func (sd *SlotDir) String() string {
	var out strings.Builder
	out.Grow(6)

	if sd.name.Val == "" {
		out.WriteString("@slot")
	} else {
		out.WriteString("@slot(")
		out.WriteString(sd.name.String())
		out.WriteString(")")
	}

	if sd.block != nil {
		out.WriteString(sd.block.String())
		out.WriteString("@end")
	}

	return out.String()
}

func (sd *SlotDir) AllChunks() []Chunk {
	if sd.block == nil {
		return []Chunk{}
	}
	return sd.block.AllChunks()
}
