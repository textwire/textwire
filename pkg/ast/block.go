package ast

import (
	"strings"

	"github.com/textwire/textwire/v4/pkg/token"
)

// Block holds chunks.
type Block struct {
	BaseNode
	Chunks []Chunk
}

func NewBlock(tok token.Token) *Block {
	return &Block{
		BaseNode: NewBaseNode(tok),
		Chunks:   make([]Chunk, 0),
	}
}

func (*Block) chunkNode() {}

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

func (b *Block) ExtractPassDirs() []*PassDir {
	passDirs := []*PassDir{}
	for _, chunk := range b.AllChunks() {
		if passDir, ok := chunk.(*PassDir); ok {
			passDirs = append(passDirs, passDir)
		}
	}
	return passDirs
}

// ToDefaultPassDir removes all *PassDir chunks and empty *Text nodes with
// whitespace from the block and returns a *PassDir with this block.
func (b *Block) ToDefaultPassDir(compName string) *PassDir {
	passDir := NewPassDir(*b.Tok(), NewStrExpr(*b.Tok(), ""))
	passDir.Block = NewBlock(*b.Tok())
	passDir.CompName = compName

	newChunks := []Chunk{}

	for _, chunk := range b.Chunks {
		switch t := chunk.(type) {
		case *PassDir:
			continue
		case *Text:
			content := strings.Trim(t.Token.Lit, " \n\t\r")
			if content == "" {
				continue
			}

			t.Token.Lit = content
			newChunks = append(newChunks, t)
		default:
			newChunks = append(newChunks, t)
		}
	}

	if len(newChunks) == 0 {
		return nil
	}

	passDir.Block.Chunks = newChunks

	return passDir
}
