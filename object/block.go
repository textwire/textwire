package object

import (
	"bytes"
	"strings"
)

type Block struct {
	Elements []Object
}

func (b *Block) Type() ObjectType {
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

	for _, e := range b.Elements {
		out.WriteString(spaces + e.Dump(ident))
	}

	return out.String()
}

func (b *Block) Val() any {
	var result []any

	for _, e := range b.Elements {
		result = append(result, e.Val())
	}

	return result
}

func (b *Block) Is(t ObjectType) bool {
	return t == b.Type()
}
