package ast

import (
	"strings"

	"github.com/textwire/textwire/v4/pkg/token"
)

type SlotDir struct {
	BaseNode
	CompName string   // Component name
	name     *StrExpr // Empty when @slot is default
}

func NewSlotDir(tok token.Token, name *StrExpr, compName string) *SlotDir {
	return &SlotDir{
		BaseNode: NewBaseNode(tok),
		name:     name,
		CompName: compName,
	}
}

func (*SlotDir) chunkNode() {}

func (sd *SlotDir) Name() *StrExpr {
	return sd.name
}

func (sd *SlotDir) String() string {
	var out strings.Builder
	out.Grow(6)

	out.WriteString(sd.Token.Lit)

	if sd.name.Val != "" {
		out.WriteString("(")
		out.WriteString(sd.name.String())
		out.WriteString(")")
	}

	return out.String()
}
