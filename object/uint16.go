package object

import "fmt"

type Uint16 struct {
	Value uint16
}

func (i *Uint16) Type() ObjectType {
	return UINT16_OBJ
}

func (i *Uint16) String() string {
	return fmt.Sprintf("%d", i.Value)
}
