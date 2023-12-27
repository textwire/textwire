package object

import "fmt"

type UnsignedInteger struct {
	Value uint64
}

func (i *UnsignedInteger) Type() ObjectType {
	return UINT_OBJ
}

func (i *UnsignedInteger) String() string {
	return fmt.Sprintf("%d", i.Value)
}
