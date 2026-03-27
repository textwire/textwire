package ast

import (
	"strings"

	"github.com/textwire/textwire/v4/pkg/token"
)

type SlotDir struct {
	BaseNode
	CompName string   // Component name
	Name     *StrExpr // Empty when @slot is default
	Block    *Block   // Optional block coming from @pass block, can be nil
}

func NewSlotDir(tok token.Token, name *StrExpr, compName string) *SlotDir {
	return &SlotDir{
		BaseNode: NewBaseNode(tok),
		Name:     name,
		CompName: compName,
	}
}

func (*SlotDir) chunkNode() {}

func (sd *SlotDir) String() string {
	var out strings.Builder
	out.Grow(6)

	out.WriteString(sd.Token.Lit)

	if sd.Name.Val != "" {
		out.WriteString("(")
		out.WriteString(sd.Name.String())
		out.WriteString(")")
	}

	return out.String()
}
