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
	return fmt.Sprintf(`<span style="%s">%d</span>`, DUMP_NUM, i.Value)
}

func (i *Int) JSON() (string, error) {
	return fmt.Sprintf("%d", i.Value), nil
}

func (i *Int) Val() any {
	return i.Value
}

func (i *Int) Is(t ObjectType) bool {
	return t == i.Type()
}
