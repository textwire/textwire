package value

import (
	"bytes"
	"strings"
)

type Block struct {
	Elements []Value
}

func (b *Block) Type() ValueType {
	return BLOCK_OBJ
}

func (b *Block) String() string {
	var out strings.Builder
	out.Grow(len(b.Elements))

	for _, e := range b.Elements {
		out.WriteString(e.String())
	}

	return out.String()
}

func (b *Block) Dump(ident int) string {
	spaces := strings.Repeat("  ", ident)
	ident += 1

	var out bytes.Buffer
	out.Grow(len(b.Elements))
	for _, e := range b.Elements {
		out.WriteString(spaces + e.Dump(ident))
	}

	return out.String()
}

func (b *Block) JSON() (string, error) {
	return "", nil
}

func (b *Block) Native() any {
	var vals []any

	for _, e := range b.Elements {
		vals = append(vals, e.Native())
	}

	return vals
}

func (b *Block) Is(t ValueType) bool {
	return t == b.Type()
}
