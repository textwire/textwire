package ast

import (
	"strings"

	"github.com/textwire/textwire/v3/pkg/token"
)

// Block holds chunks.
type Block struct {
	BaseNode
	Chunks []Chunk
}

func NewBlock(tok token.Token) *Block {
	return &Block{
		BaseNode: NewBaseNode(tok),
	}
}

func (_ *Block) chunkNode() {}

func (b *Block) String() string {
	var out strings.Builder
	out.Grow(len(b.Chunks))

	for i := range b.Chunks {
		out.WriteString(b.Chunks[i].String())
	}

	return out.String()
}

func (b *Block) AllChunks() []Chunk {
	if b.Chunks == nil {
		return []Chunk{}
	}

	chunks := make([]Chunk, 0, len(b.Chunks))

	for _, chunk := range b.Chunks {
		if chunk == nil {
			continue
		}

		if s, ok := chunk.(NodeWithChunks); ok {
			chunks = append(chunks, s.(Chunk))
			chunks = append(chunks, s.AllChunks()...)
		}
	}

	return chunks
}
