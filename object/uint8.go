package object

import "fmt"

type Uint8 struct {
	Value uint8
}

func (i *Uint8) Type() ObjectType {
	return UINT8_OBJ
}

func (i *Uint8) String() string {
	return fmt.Sprintf("%d", i.Value)
}
