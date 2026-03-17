package ast

import (
	"fmt"
	"strings"

	"github.com/textwire/textwire/v4/pkg/token"
)

type IfDir struct {
	BaseNode
	Cond       Expression
	IfBlock    *Block // @if()<IfBlock>@end
	ElseBlock  *Block // @else<ElseBlock>@end
	ElseifDirs []*ElseIfDir
}

func NewIfDir(tok token.Token) *IfDir {
	return &IfDir{
		BaseNode: NewBaseNode(tok),
	}
}

func (*IfDir) chunkNode() {}

func (id *IfDir) String() string {
	var out strings.Builder
	out.Grow(20 + len(id.ElseifDirs)*2)

	fmt.Fprintf(&out, "@if(%s)\n%s", id.Cond, id.IfBlock)

	for _, e := range id.ElseifDirs {
		out.WriteString(e.String())
	}

	if id.ElseBlock != nil {
		out.WriteString("@else\n")
		out.WriteString(id.ElseBlock.String() + "\n")
	}

	out.WriteString("@end\n")

	return out.String()
}

func (id *IfDir) AllChunks() []Chunk {
	chunks := make([]Chunk, 0)
	if id.IfBlock != nil {
		chunks = append(chunks, id.IfBlock.AllChunks()...)
	}

	if id.ElseBlock != nil {
		chunks = append(chunks, id.ElseBlock.AllChunks()...)
	}

	for _, dir := range id.ElseifDirs {
		chunks = append(chunks, dir.AllChunks()...)
	}

	return chunks
}
