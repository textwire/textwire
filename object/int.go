package object

import (
	"fmt"
)

type Int struct {
	Value int64
}

func (i *Int) Type() ObjectType {
	return INT_OBJ
}

func (i *Int) String() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Int) Dump(ident int) string {
	return fmt.Sprintf("<span class='textwire-num'>%d</span>", i.Value)
}

func (i *Int) Val() interface{} {
	return i.Value
}

func (i *Int) Is(t ObjectType) bool {
	return t == i.Type()
}
