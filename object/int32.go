package object

import "fmt"

type Int32 struct {
	Value int32
}

func (i *Int32) Type() ObjectType {
	return INT32_OBJ
}

func (i *Int32) String() string {
	return fmt.Sprintf("%d", i.Value)
}
