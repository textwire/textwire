package ast

import (
	"strings"

	"github.com/textwire/textwire/v4/pkg/token"
)

type ProvideDir struct {
	BaseNode
	CompName string     // Component name
	Name     *StrExpr   // Cannot be empty
	Block    *Block     // Optional block statement, can be nil
	Cond     Expression // When you have @provideif, this field will be boolean expression
}

func NewProvideDir(tok token.Token, name *StrExpr) *ProvideDir {
	return &ProvideDir{
		BaseNode: NewBaseNode(tok),
		Name:     name,
	}
}

func (*ProvideDir) chunkNode() {}

func (pd *ProvideDir) String() string {
	var out strings.Builder
	out.Grow(6)

	out.WriteString(pd.Token.Lit)

	hasParens := pd.Cond != nil || pd.Name.Val != ""

	if hasParens {
		out.WriteString("(")
	}

	if pd.Cond != nil {
		out.WriteString(pd.Cond.String())

		if pd.Name.Val != "" {
			out.WriteString(", ")
		}
	}

	out.WriteString(pd.Name.String())

	if hasParens {
		out.WriteString(")")
	}

	if pd.Block != nil {
		out.WriteString(pd.Block.String())
		out.WriteString("@end")
	}

	return out.String()
}

func (pd *ProvideDir) AllChunks() []Chunk {
	if pd.Block == nil {
		return []Chunk{}
	}
	return pd.Block.AllChunks()
}
