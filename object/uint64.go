package object

import "fmt"

type Uint64 struct {
	Value uint64
}

func (i *Uint64) Type() ObjectType {
	return UINT64_OBJ
}

func (i *Uint64) String() string {
	return fmt.Sprintf("%d", i.Value)
}
