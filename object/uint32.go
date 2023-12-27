package object

import "fmt"

type Uint32 struct {
	Value uint32
}

func (i *Uint32) Type() ObjectType {
	return UINT32_OBJ
}

func (i *Uint32) String() string {
	return fmt.Sprintf("%d", i.Value)
}
