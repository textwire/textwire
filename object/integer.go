package object

import "fmt"

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

func (i *Integer) String() string {
	return fmt.Sprintf("%d", i.Value)
}
