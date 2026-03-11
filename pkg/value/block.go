package value

import (
	"strings"
)

type Block struct {
	Elements []Value
}

func (*Block) Type() ValueType {
	return BLOCK_VAL
}

func (b *Block) String() string {
	var out strings.Builder
	out.Grow(len(b.Elements))

	for _, e := range b.Elements {
		out.WriteString(e.String())
	}

	return out.String()
}

func (b *Block) Is(t ValueType) bool {
	return t == b.Type()
}
