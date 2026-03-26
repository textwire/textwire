package ast

import (
	"strings"

	"github.com/textwire/textwire/v4/pkg/token"
)

type CompDir struct {
	BaseNode
	Name     *StrExpr // Relative path to the component 'components/book'
	Argument *ObjExpr
	CompProg *Program // AST node of the component file Name
	Block    *Block
}

func NewCompDir(tok token.Token) *CompDir {
	return &CompDir{
		BaseNode: NewBaseNode(tok),
	}
}

func (*CompDir) chunkNode() {}

func (cd *CompDir) String() string {
	var out strings.Builder
	out.Grow(6)

	out.WriteString(cd.Block.Token.Lit)
	out.WriteByte('(')

	out.WriteString(cd.Name.Val)

	if cd.Argument != nil {
		out.WriteString(", ")
		out.WriteString(cd.Argument.String())
	}

	out.WriteByte(')')
	out.WriteString(cd.Block.String())
	out.WriteString("@end")

	return out.String()
}

func (cd *CompDir) AllChunks() []Chunk {
	if cd.CompProg == nil {
		return []Chunk{}
	}

	return cd.CompProg.AllChunks()
}
