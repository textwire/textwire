package ast

import (
	"strings"

	"github.com/textwire/textwire/v4/pkg/token"
)

type PassDir struct {
	BaseNode
	CompName string     // Component name
	Name     *StrExpr   // Cannot be empty
	Block    *Block     // Optional block statement, can be nil
	Cond     Expression // When you have @passif, this field will be boolean expression
}

func NewPassDir(tok token.Token, name *StrExpr) *PassDir {
	return &PassDir{
		BaseNode: NewBaseNode(tok),
		Name:     name,
	}
}

func (*PassDir) chunkNode() {}

func (pd *PassDir) String() string {
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

func (pd *PassDir) AllChunks() []Chunk {
	if pd.Block == nil {
		return []Chunk{}
	}
	return pd.Block.AllChunks()
}

func (pd *PassDir) IsDefault() bool {
	return pd.Name.Val == ""
}
