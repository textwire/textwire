package object

import "fmt"

type Int8 struct {
	Value int8
}

func (i *Int8) Type() ObjectType {
	return INT8_OBJ
}

func (i *Int8) String() string {
	return fmt.Sprintf("%d", i.Value)
}
