package ast

import (
	"fmt"
	"strings"

	"github.com/textwire/textwire/v4/pkg/token"
)

type ForDir struct {
	BaseNode
	Init      Statement  // Initialization statement
	Cond      Expression // Condition expression
	Post      Statement  // Post iteration statement
	ElseBlock *Block     // @else block
	Block     *Block
}

func NewForDir(tok token.Token) *ForDir {
	return &ForDir{
		BaseNode: NewBaseNode(tok),
	}
}

func (*ForDir) chunkNode() {}

func (fd *ForDir) LoopBlock() *Block {
	if fd.Block == nil {
		panic("Block must not be nil on ForStmt when calling LoopBlock()")
	}
	return fd.Block
}

func (fd *ForDir) String() string {
	var out strings.Builder
	out.Grow(20)

	fmt.Fprintf(&out, "@for(%s; %s; %s)\n", fd.Init, fd.Cond, fd.Post)

	out.WriteString(fd.Block.String() + "\n")

	if fd.ElseBlock != nil {
		out.WriteString("@else\n")
		out.WriteString(fd.ElseBlock.String() + "\n")
	}

	out.WriteString("@end\n")

	return out.String()
}

func (fd *ForDir) AllChunks() []Chunk {
	chunks := make([]Chunk, 0)
	if fd.Block != nil {
		chunks = append(chunks, fd.Block.AllChunks()...)
	}

	if fd.ElseBlock != nil {
		chunks = append(chunks, fd.ElseBlock.AllChunks()...)
	}

	return chunks
}
