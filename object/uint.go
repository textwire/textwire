package object

import "fmt"

type Uint struct {
	Value uint
}

func (i *Uint) Type() ObjectType {
	return UINT_OBJ
}

func (i *Uint) String() string {
	return fmt.Sprintf("%d", i.Value)
}
