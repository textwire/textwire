package ast

import (
	"strings"

	"github.com/textwire/textwire/v4/pkg/token"
)

type CompDir struct {
	BaseNode
	Name        *StrExpr // Relative path to the component 'components/book'
	Argument    *ObjExpr
	CompProg    *Program // AST node of the component file Name
	DefaultPass *PassDir // Default pass with no name, can be nil
	Passes      []*PassDir
}

func NewCompDir(tok token.Token) *CompDir {
	return &CompDir{
		BaseNode: NewBaseNode(tok),
		Passes:   make([]*PassDir, 0),
	}
}

func (*CompDir) chunkNode() {}

func (cd *CompDir) String() string {
	var out strings.Builder
	out.Grow(6)

	out.WriteString(cd.Token.Lit)
	out.WriteByte('(')

	out.WriteString(cd.Name.String())

	if cd.Argument != nil {
		out.WriteString(", ")
		out.WriteString(cd.Argument.String())
	}

	out.WriteByte(')')

	if cd.DefaultPass != nil {
		out.WriteString(cd.DefaultPass.String())
	}

	for i := range cd.Passes {
		out.WriteString(cd.Passes[i].String())
	}

	out.WriteString("@end")

	return out.String()
}

func (cd *CompDir) AllChunks() []Chunk {
	if cd.CompProg == nil {
		return []Chunk{}
	}

	return cd.CompProg.AllChunks()
}
