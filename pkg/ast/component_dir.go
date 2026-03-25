package ast

import (
	"strings"

	"github.com/textwire/textwire/v4/pkg/token"
)

type CompDir struct {
	BaseNode
	Name     *StrExpr // Relative path to the component 'components/book'
	Argument *ObjExpr
	CompProg *Program      // AST node of the component file Name
	Provides []*ProvideDir // Each slot of the component's block
}

func NewCompDir(tok token.Token) *CompDir {
	return &CompDir{
		BaseNode: NewBaseNode(tok),
	}
}

func (*CompDir) chunkNode() {}

func (cd *CompDir) ArgsString() string {
	var out strings.Builder
	out.Grow(10)

	out.WriteString(cd.Name.String())

	if cd.Argument != nil {
		out.WriteString(", ")
		out.WriteString(cd.Argument.String())
	}

	return out.String()
}

func (cd *CompDir) String() string {
	var out strings.Builder
	out.Grow(len(cd.Provides) + 20)

	out.WriteString("@component(")
	out.WriteString(cd.ArgsString())
	out.WriteString(")")

	for _, slot := range cd.Provides {
		out.WriteString(slot.String())
	}

	if len(cd.Provides) > 0 {
		out.WriteString("@end")
	}

	return out.String()
}

func (cd *CompDir) AllChunks() []Chunk {
	if cd.CompProg == nil {
		return []Chunk{}
	}

	return cd.CompProg.AllChunks()
}
