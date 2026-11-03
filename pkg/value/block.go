package value

import (
	"strings"
)

type Block struct {
	Chunks []Value

	// hasBreak and hasContinue cache whether a break or continue is anywhere
	// in the block, so checking them doesn't walk nested blocks.
	hasBreak    bool
	hasContinue bool
}

func NewBlock(cap int) *Block {
	return &Block{Chunks: make([]Value, 0, cap)}
}

// Append adds a chunk to the block and updates the cached control flags.
func (b *Block) Append(chunk Value) {
	b.Chunks = append(b.Chunks, chunk)

	inner, isBlock := chunk.(*Block)

	if chunk.Is(BREAK_VAL) || (isBlock && inner.hasBreak) {
		b.hasBreak = true
	}

	if chunk.Is(CONTINUE_VAL) || (isBlock && inner.hasContinue) {
		b.hasContinue = true
	}
}

// HasControl reports whether the block contains a break or continue value of
// the given type, including inside nested blocks.
func (b *Block) HasControl(t ValueType) bool {
	switch t {
	case BREAK_VAL:
		return b.hasBreak
	case CONTINUE_VAL:
		return b.hasContinue
	}

	return false
}

func (*Block) Type() ValueType {
	return BLOCK_VAL
}

func (b *Block) String() string {
	var out strings.Builder
	out.Grow(len(b.Chunks))

	for i := range b.Chunks {
		out.WriteString(b.Chunks[i].String())
	}

	return out.String()
}

func (b *Block) Is(t ValueType) bool {
	return t == b.Type()
}
