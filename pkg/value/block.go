package value

import (
	"strings"
)

type Block struct {
	Chunks []Value
}

func NewBlock(cap int) *Block {
	return &Block{Chunks: make([]Value, 0, cap)}
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
