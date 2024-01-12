package object

import "bytes"

type Block struct {
	Elements []Object
}

func (b *Block) Type() ObjectType {
	return BLOCK_OBJ
}

func (b *Block) String() string {
	var result bytes.Buffer

	for _, e := range b.Elements {
		result.WriteString(e.String())
	}

	return result.String()
}
