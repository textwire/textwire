package ast

import (
	"fmt"

	"github.com/textwire/textwire/v3/pkg/token"
)

type UseDir struct {
	BaseNode
	Name       *StrExpr              // Relative path to the layout like 'layouts/main'
	LayoutProg *Program              // AST node of the layout file Name
	Inserts    map[string]*InsertDir // @use connection to @insert directives
}

func NewUseDir(tok token.Token) *UseDir {
	return &UseDir{
		BaseNode: NewBaseNode(tok),
	}
}

func (_ *UseDir) chunkNode() {}

func (ud *UseDir) String() string {
	return fmt.Sprintf(`@use(%s)`, ud.Name)
}

func (ud *UseDir) AllChunks() []Chunk {
	if ud.LayoutProg == nil {
		return []Chunk{}
	}
	return ud.LayoutProg.AllChunks()
}
