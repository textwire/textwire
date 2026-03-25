package ast

import (
	"strings"

	"github.com/textwire/textwire/v4/pkg/token"
)

// ProvideifDir cannot be external and live in component file.
type ProvideifDir struct {
	BaseNode  // @provideif(bool, 'name'?)
	Cond      Expression
	CompName  string   // Component name
	isDefault bool     // Whether the slot is named or default
	name      *StrExpr // Empty when @slot is default
	block     *Block   // Required, cannot be nil
}

func NewProvideifDir(tok token.Token, name *StrExpr, compName string) *ProvideifDir {
	return &ProvideifDir{
		BaseNode: NewBaseNode(tok),
		CompName: compName,
		name:     name,
	}
}

func (*ProvideifDir) chunkNode() {}

func (sd *ProvideifDir) Name() *StrExpr {
	return sd.name
}

func (sd *ProvideifDir) IsDefault() bool {
	return sd.isDefault
}

func (sd *ProvideifDir) SetIsDefault(val bool) {
	sd.isDefault = val
}

func (sd *ProvideifDir) SetBlock(b *Block) {
	sd.block = b
}

func (sd *ProvideifDir) Block() *Block {
	return sd.block
}

func (sd *ProvideifDir) String() string {
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

func (sd *ProvideifDir) AllChunks() []Chunk {
	if sd.block == nil {
		return []Chunk{}
	}
	return sd.block.AllChunks()
}
