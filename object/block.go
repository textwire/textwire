package object

import (
	"bytes"
)

type Block struct {
	Elements []Object
}

func (b *Block) Type() ObjectType {
	return BLOCK_OBJ
}

func (b *Block) String() string {
	var out bytes.Buffer

	for _, e := range b.Elements {
		out.WriteString(e.String())
	}

	return out.String()
}

func (b *Block) Val() interface{} {
	var result []interface{}

	for _, e := range b.Elements {
		result = append(result, e.Val())
	}

	return result
}

func (b *Block) Is(t ObjectType) bool {
	return t == b.Type()
}
