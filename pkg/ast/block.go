package ast

import (
	"strings"

	"codeberg.org/textwire/textwire/v6/pkg/token"
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

// ToDefaultPassDir removes all *PassDir chunks from the block and returns a
// *PassDir with the remaining chunks.
func (b *Block) ToDefaultPassDir(compName string) *PassDir {
	passDir := NewPassDir(*b.Tok(), NewStrExpr(*b.Tok(), ""))
	passDir.Block = NewBlock(*b.Tok())
	passDir.CompName = compName

	newChunks := []Chunk{}

	hasContent := false

	for _, chunk := range b.Chunks {
		if _, ok := chunk.(*PassDir); ok {
			continue
		}

		newChunks = append(newChunks, chunk)
		if !isBlankText(chunk) {
			hasContent = true
		}
	}

	// Blank text and @pass alone don't make a default slot, so the component
	// keeps its default pass nil.
	if !hasContent {
		return nil
	}

	passDir.Block.Chunks = newChunks

	return passDir
}

// isBlankText reports whether chunk is a text node with only whitespace.
func isBlankText(chunk Chunk) bool {
	text, ok := chunk.(*Text)
	if !ok {
		return false
	}
	return strings.TrimSpace(text.Token.Lit) == ""
}
