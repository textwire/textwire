package value

import (
	"slices"
	"strings"
)

type Block struct {
	chunks []Value

	// hasBreak and hasContinue cache whether a break or continue is anywhere
	// in the block, so checking them doesn't walk nested blocks.
	hasBreak    bool
	hasContinue bool
}

func NewBlock(cap int) *Block {
	return &Block{chunks: make([]Value, 0, cap)}
}

func (b *Block) Chunks() []Value {
	return slices.Clone(b.chunks)
}

// Append adds a chunk to the block and updates the cached control flags.
func (b *Block) Append(chunk Value) {
	b.chunks = append(b.chunks, chunk)

	inner, isBlock := chunk.(*Block)

	if chunk.Is(BREAK_VAL) || (isBlock && inner.hasBreak) {
		b.hasBreak = true
	}

	if chunk.Is(CONTINUE_VAL) || (isBlock && inner.hasContinue) {
		b.hasContinue = true
	}
}

func (b *Block) HasBreak() bool {
	return b.hasBreak
}

func (b *Block) HasContinue() bool {
	return b.hasContinue
}

func (*Block) Type() ValueType {
	return BLOCK_VAL
}

func (b *Block) String() string {
	var out strings.Builder
	out.Grow(len(b.chunks))

	for i := range b.chunks {
		out.WriteString(b.chunks[i].String())
	}

	return out.String()
}

func (b *Block) Is(t ValueType) bool {
	return t == b.Type()
}
