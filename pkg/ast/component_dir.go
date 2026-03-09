package ast

import (
	"strings"

	"github.com/textwire/textwire/v3/pkg/token"
)

type ComponentDir struct {
	BaseNode
	Name     *StrExpr // Relative path to the component 'components/book'
	Argument *ObjExpr
	CompProg *Program        // AST node of the component file Name
	Slots    []SlotDirective // Each slot of the component's block
}

func NewComponentDir(tok token.Token) *ComponentDir {
	return &ComponentDir{
		BaseNode: NewBaseNode(tok),
	}
}

func (_ *ComponentDir) chunkNode() {}

func (_ *ComponentDir) Kind() ChunkKind {
	return ChunkKindDirective
}

func (cd *ComponentDir) ArgsString() string {
	var out strings.Builder
	out.Grow(10)

	out.WriteString(cd.Name.String())

	if cd.Argument != nil {
		out.WriteString(", ")
		out.WriteString(cd.Argument.String())
	}

	return out.String()
}

func (cd *ComponentDir) String() string {
	var out strings.Builder
	out.Grow(len(cd.Slots) + 20)

	out.WriteString("@component(")
	out.WriteString(cd.ArgsString())
	out.WriteString(")")

	for _, slot := range cd.Slots {
		out.WriteString("\n")
		out.WriteString(slot.String())
	}

	if len(cd.Slots) > 0 {
		out.WriteString("\n@end\n")
	}

	return out.String()
}

func (cd *ComponentDir) AllChunks() []Chunk {
	if cd.CompProg == nil {
		return []Chunk{}
	}

	return cd.CompProg.AllChunks()
}
