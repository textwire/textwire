package ast

import (
	"strings"

	"github.com/textwire/textwire/v4/pkg/token"
)

type ProvideDir struct {
	BaseNode
	CompName  string   // Component name
	isDefault bool     // Whether the slot is named or default
	name      *StrExpr // Empty when @slot is default
	block     *Block   // Optional block statement, can be nil
}

func NewProvideDir(tok token.Token, name *StrExpr, compName string) *ProvideDir {
	return &ProvideDir{
		BaseNode: NewBaseNode(tok),
		name:     name,
		CompName: compName,
	}
}

func (*ProvideDir) chunkNode() {}

func (pd *ProvideDir) Name() *StrExpr {
	return pd.name
}

func (pd *ProvideDir) IsDefault() bool {
	return pd.isDefault
}

func (pd *ProvideDir) SetIsDefault(val bool) {
	pd.isDefault = val
}

func (pd *ProvideDir) SetBlock(b *Block) {
	pd.block = b
}

func (pd *ProvideDir) Block() *Block {
	return pd.block
}

func (pd *ProvideDir) String() string {
	var out strings.Builder
	out.Grow(6)

	out.WriteString(pd.Token.Lit)

	if pd.name.Val != "" {
		out.WriteString("(")
		out.WriteString(pd.name.String())
		out.WriteString(")")
	}

	if pd.block != nil {
		out.WriteString(pd.block.String())
		out.WriteString("@end")
	}

	return out.String()
}

func (pd *ProvideDir) AllChunks() []Chunk {
	if pd.block == nil {
		return []Chunk{}
	}
	return pd.block.AllChunks()
}
