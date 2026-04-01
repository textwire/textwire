package ast

import (
	"fmt"
	"strings"

	"github.com/textwire/textwire/v4/pkg/token"
)

type EachDir struct {
	BaseNode
	Var       *IdentExpr // Variable name
	Arr       Expression // Arr to loop over
	ElseBlock *Block     // @else<ElseBlock>@end
	Block     *Block
}

func NewEachDir(tok token.Token) *EachDir {
	return &EachDir{
		BaseNode: NewBaseNode(tok),
	}
}

func (*EachDir) chunkNode() {}

func (ed *EachDir) LoopBlock() *Block {
	if ed.Block == nil {
		panic("Block must not be nil on EachStmt when calling LoopBlock()")
	}
	return ed.Block
}

func (ed *EachDir) String() string {
	var out strings.Builder
	out.Grow(26)

	fmt.Fprintf(&out, "@each(%s in %s)\n%s\n", ed.Var, ed.Arr, ed.Block)

	if ed.ElseBlock != nil {
		out.WriteString("@else\n")
		out.WriteString(ed.ElseBlock.String() + "\n")
	}

	out.WriteString("@end\n")

	return out.String()
}

func (ed *EachDir) AllChunks() []Chunk {
	chunks := make([]Chunk, 0)
	if ed.Block != nil {
		chunks = append(chunks, ed.Block.AllChunks()...)
	}

	if ed.ElseBlock != nil {
		chunks = append(chunks, ed.ElseBlock.AllChunks()...)
	}

	return chunks
}
