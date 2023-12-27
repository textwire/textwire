package object

import "fmt"

type Int struct {
	Value int
}

func (i *Int) Type() ObjectType {
	return INT_OBJ
}

func (i *Int) String() string {
	return fmt.Sprintf("%d", i.Value)
}
