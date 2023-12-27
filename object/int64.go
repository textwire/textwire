package object

import "fmt"

type Int64 struct {
	Value int64
}

func (i *Int64) Type() ObjectType {
	return INT64_OBJ
}

func (i *Int64) String() string {
	return fmt.Sprintf("%d", i.Value)
}
