package ast

import (
	"fmt"
	"strings"

	"github.com/textwire/textwire/v4/pkg/token"
)

type InsertDir struct {
	BaseNode
	Name *StrExpr
	// Argument is nil when Block is present
	Argument Expression
	// Block is nil when Argument is present
	Block *Block
	// AbsPath of the file where insert is located
	AbsPath string
}

func NewInsertDir(tok token.Token, absPath string) *InsertDir {
	return &InsertDir{
		BaseNode: NewBaseNode(tok),
		AbsPath:  absPath,
	}
}

func (*InsertDir) chunkNode() {}

func (i *InsertDir) String() string {
	var out strings.Builder
	out.Grow(30)

	if i.Argument != nil {
		fmt.Fprintf(&out, `@insert("%s", %s)`, i.Name, i.Argument)
		return out.String()
	}

	fmt.Fprintf(&out, `@insert("%s")`, i.Name)
	out.WriteString(i.Block.String())
	out.WriteString(`@end`)

	return out.String()
}

func (i *InsertDir) AllChunks() []Chunk {
	if i.Block == nil {
		return []Chunk{}
	}
	return i.Block.AllChunks()
}
