package object

import "fmt"

type Int16 struct {
	Value int16
}

func (i *Int16) Type() ObjectType {
	return INT16_OBJ
}

func (i *Int16) String() string {
	return fmt.Sprintf("%d", i.Value)
}
